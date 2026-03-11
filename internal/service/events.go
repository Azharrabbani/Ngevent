package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"ngevent/internal/repository"
	"ngevent/internal/utils"
	"ngevent/internal/utils/helper"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type EventService struct {
	EventRepo          repository.EventsRepo
	UserRepo           repository.UsersRepo
	ProfileRepo        repository.OrganizerProfileRepo
	CategoryRepo       repository.CategoriesRepo
	EmailTaskPublisher NewTaskEmail
	rdb                *redis.Client
}

func NewEventService(
	eventRepo repository.EventsRepo,
	userRepo repository.UsersRepo,
	profileRepo repository.OrganizerProfileRepo,
	categoryRepo repository.CategoriesRepo,
	emailTaskPublisher NewTaskEmail,
	rdb *redis.Client,
) *EventService {
	return &EventService{
		EventRepo:    eventRepo,
		UserRepo:     userRepo,
		ProfileRepo:  profileRepo,
		CategoryRepo: categoryRepo,
		rdb:          rdb,
	}
}

var (
	eventBannerPath = "./storage/event/banner"
)

func (s *EventService) InvalidateCache() {
	ctx := context.Background()

	patterns := []string{
		"events:all:*",
	}

	for _, pattern := range patterns {
		iter := s.rdb.Scan(ctx, 0, pattern, 0).Iterator()
		for iter.Next(ctx) {
			s.rdb.Del(ctx, iter.Val())
		}
	}

	// Use SCAN for pattern keys to avoid blocking
	log.Println("[CACHE] events cache invalidated")
}

func (s *EventService) CreateEvent(req *dto.EventReq) error {
	// Search the eo profile
	profile, err := s.ProfileRepo.FindByUserID(req.UserID)
	if err != nil {
		return errors.New("profile not found")
	}

	// Validate len categories & tickets
	if len(req.Categories) < 0 || len(req.Tickets) < 0 {
		return errors.New("categories or ticket cannot be empty")
	}

	// Get categories
	categories, err := s.CategoryRepo.FindByIDs(req.Categories)
	if err != nil {
		return err
	}

	// Declared the tickets
	var tickets []*model.Tickets
	for _, ticket := range req.Tickets {
		tickets = append(tickets, &model.Tickets{
			Name:       ticket.Name,
			Price:      ticket.Price,
			Quantity:   ticket.Quantity,
			TicketType: ticket.TicketType,
		})
	}

	// Convert unix to date
	date := helper.ConvertUnixtoDate(req.Date)

	// Default status
	status := "pending"
	if req.Status != nil {
		status = *req.Status
	}

	// Get coordinates from req
	coordinates := fmt.Sprintf("POINT(%f %f)", req.Address.Long, req.Address.Lat)

	// Save events
	event := &model.Events{
		ProfileID:     profile.ID,
		Name:          req.Name,
		Slug:          utils.CreateSlug(req.Name),
		Status:        status,
		Description:   req.Description,
		Address:       req.Address.Address,
		City:          req.Address.City,
		Country:       req.Address.Country,
		DetailAddress: req.Address.DetailAddress,
		Coordinates:   coordinates,
		Date:          date,
	}

	if err := s.EventRepo.Create(event, categories, tickets); err != nil {
		return err
	}

	// Invalidate cache after update
	s.InvalidateCache()

	// Email the admins
	admins, err := s.UserRepo.FindByRole("admin")
	if err != nil {
		log.Println("[ERROR] admin data not found")
	}

	for _, admin := range admins {
		AdminEmailPayload := &model.EventEmailPayload{
			To:        admin.Email,
			EOName:    profile.Name,
			EOEmail:   profile.User.Email,
			EventName: event.Name,
		}
		if err := s.EmailTaskPublisher.Enqueue(model.TypeEventAdminNotification, AdminEmailPayload); err != nil {
			log.Printf("[EMAIL] failed sending email to admin %s\n", admin.Email)
		}
	}

	// Email the EO
	EOEmailPayload := &model.EventEmailPayload{
		To:        profile.User.Email,
		EventName: event.Name,
	}
	if err := s.EmailTaskPublisher.Enqueue(model.TypeEventEONotification, EOEmailPayload); err != nil {
		log.Printf("[EMAIL] failed sending email to event organizer %s\n", profile.User.Email)
	}

	return nil
}

func (s *EventService) ReviewEvent(req *dto.ReviewEventReq) error {
	// Validate the event
	event, err := s.EventRepo.FindByID(req.ID)
	if err != nil {
		return errors.New("event not found")
	}

	event.Status = req.Status

	if err := s.EventRepo.ReviewEvent(event); err != nil {
		return err
	}

	// Invalidate cache after update
	s.InvalidateCache()

	// Email the EO
	organizer, err := s.ProfileRepo.FindByID(event.ProfileID)
	if err != nil {
		log.Println("[ERROR] organizer data not found")
	}

	EOEmailPayload := &model.EventEmailPayload{
		To:        organizer.User.Email,
		EventName: event.Name,
	}
	if err := s.EmailTaskPublisher.Enqueue(model.TypeEventEONotification, EOEmailPayload); err != nil {
		log.Printf("[EMAIL] failed sending email to event organizer %s\n", organizer.User.Email)
	}

	return nil
}

func (s *EventService) GetEvents(pagination model.Pagination) (*model.PaginationRow[*dto.EventsResp], error) {
	return s.EventRepo.FindAll(pagination)
}

func (s *EventService) GeteventsBySlug(title string, pagination model.Pagination) (*model.PaginationRow[*dto.EventsResp], error) {
	slugTitle := utils.CreateSlug(title)

	return s.EventRepo.FindBySlug(slugTitle, pagination)
}

func (s *EventService) GetEventByID(id string) (*dto.EventsResp, error) {
	// Search the event
	event, err := s.EventRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("event not found")
	}

	// Get the organizer profile
	organizer, err := s.ProfileRepo.FindByID(event.ID)
	if err != nil {
		return nil, errors.New("organizer not found")
	}

	var eventCategories []dto.EventCategories
	for _, category := range event.Categories {
		eventCategories = append(eventCategories, dto.EventCategories{
			ID:   category.ID,
			Name: category.Category.Name,
		})
	}

	var tickets []dto.Tickets
	for _, ticket := range event.Tickets {
		tickets = append(tickets, dto.Tickets{
			ID:         ticket.ID,
			Name:       ticket.Name,
			Price:      ticket.Price,
			Quantity:   ticket.Quantity,
			TicketType: ticket.TicketType,
		})
	}

	respReq := &dto.EventRespReq{
		Event:           event,
		Organizer:       organizer,
		EventCategories: eventCategories,
		Tickets:         tickets,
		Date:            helper.ConvertDatetoUnix(event.Date.Format(time.RFC3339)),
		CreatedAt:       helper.ConvertDatetoUnix(event.CreatedAt.Format(time.RFC3339)),
		UpdatedAt:       helper.ConvertDatetoUnix(event.UpdatedAt.Format(time.RFC3339)),
		DeletedAt:       helper.ConvertDatetoUnix(event.DeletedAt.Format(time.RFC3339)),
	}
	eventResp, err := dto.ToEventResp(respReq)
	if err != nil {
		return nil, err
	}

	return eventResp, nil
}

func (s *EventService) UpdateEventBanner(req *dto.UpdateEventReq) error {
	// Validate event
	event, err := s.EventRepo.FindByID(req.ID)
	if err != nil {
		return errors.New("event not found")
	}

	organizer, err := s.ProfileRepo.FindByUserID(req.UserID)
	if err != nil {
		return errors.New("organizer not found")
	}

	// Validate user
	if event.ProfileID != organizer.ID {
		return errors.New("unauthorized action")
	}

	// Validate image
	if err := helper.ValidateImage(req.Banner); err != nil {
		return err
	}

	// Save to local
	_, fileName, err := helper.SaveToLocal(req.Banner, eventBannerPath)
	if err != nil {
		return err
	}

	// Update event banner
	if err := s.EventRepo.UpdateBannerEvent(event.ID, fileName); err != nil {
		log.Printf("[ERROR] Failed update the banner %v\n", err)
		return err
	}

	// Remove old banner if exist
	if event.Banner != nil {
		oldBanner := fmt.Sprintf("%s/%s", eventBannerPath, *event.Banner)
		if err := os.Remove(oldBanner); err != nil {
			log.Printf("failed to remove file from local %v\n", err)
		}
	}

	// Invalidate cache after update
	s.InvalidateCache()

	return nil
}

// func (s *EventService) UpdateEvent(req *dto.EventReq) error {
// 	// Search the eo profile
// 	profile, err := s.ProfileRepo.FindByUserID(req.UserID)
// 	if err != nil {
// 		return errors.New("profile not found")
// 	}

// 	// Validate event
// 	event, err := s.EventRepo.FindByID(req.EventID)
// 	if err != nil {
// 		return errors.New("event not found")
// 	}

// 	// Validate the user
// 	if event.ProfileID != profile.ID {
// 		return errors.New("unauthorized action")
// 	}

// 	// Validate len categories & tickets
// 	if len(req.Categories) < 0 || len(req.Tickets) < 0 {
// 		return errors.New("categories or ticket cannot be empty")
// 	}

// 	// Get categories
// 	categories, err := s.CategoryRepo.FindByIDs(req.Categories)
// 	if err != nil {
// 		return err
// 	}

// 	// Declared the tickets
// 	var tickets []*model.Tickets
// 	for _, ticket := range req.Tickets {
// 		tickets = append(tickets, &model.Tickets{
// 			Name:       ticket.Name,
// 			Price:      ticket.Price,
// 			Quantity:   ticket.Quantity,
// 			TicketType: ticket.TicketType,
// 		})
// 	}

// 	// Convert unix to date
// 	date := helper.ConvertUnixtoDate(req.Date)

// 	// Get coordinates from req
// 	coordinates := fmt.Sprintf("POINT(%f %f)", req.Address.Long, req.Address.Lat)

	
// 	// Save events
// 	event = &model.Events{
// 		ProfileID:     profile.ID,
// 		Name:          req.Name,
// 		Slug:          utils.CreateSlug(req.Name),
// 		Status:        status,
// 		Description:   req.Description,
// 		Address:       req.Address.Address,
// 		City:          req.Address.City,
// 		Country:       req.Address.Country,
// 		DetailAddress: req.Address.DetailAddress,
// 		Coordinates:   coordinates,
// 		Date:          date,
// 	}

// 	if err := s.EventRepo.Create(event, categories, tickets); err != nil {
// 		return err
// 	}

// 	// Invalidate cache after update
// 	s.InvalidateCache()

// 	// Email the admins
// 	admins, err := s.UserRepo.FindByRole("admin")
// 	if err != nil {
// 		log.Println("[ERROR] admin data not found")
// 	}

// 	for _, admin := range admins {
// 		AdminEmailPayload := &model.EventEmailPayload{
// 			To:        admin.Email,
// 			EOName:    profile.Name,
// 			EOEmail:   profile.User.Email,
// 			EventName: event.Name,
// 		}
// 		if err := s.EmailTaskPublisher.Enqueue(model.TypeEventAdminNotification, AdminEmailPayload); err != nil {
// 			log.Printf("[EMAIL] failed sending email to admin %s\n", admin.Email)
// 		}
// 	}

// 	// Email the EO
// 	EOEmailPayload := &model.EventEmailPayload{
// 		To:        profile.User.Email,
// 		EventName: event.Name,
// 	}
// 	if err := s.EmailTaskPublisher.Enqueue(model.TypeEventEONotification, EOEmailPayload); err != nil {
// 		log.Printf("[EMAIL] failed sending email to event organizer %s\n", profile.User.Email)
// 	}

// 	return nil
// }
