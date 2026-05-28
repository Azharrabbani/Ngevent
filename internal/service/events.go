package service

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"ngevent/internal/repository"
	"ngevent/internal/utils"
	"ngevent/internal/utils/analytics"
	"ngevent/internal/utils/helper"
	"os"
	"path/filepath"
	"time"

	"github.com/redis/go-redis/v9"
)

type NewTaskEventExpiryPublisher interface {
	EnqueueEventExpiry(payload *model.EventExpiredPayload, endTime time.Time) error
	EnqueueUpdatedEventExpiry(payload *model.UpdatedEventExpiredPayload, endTime time.Time) error
	CancelEventExpiry(eventID string) error
}

type EventService struct {
	EventRepo            repository.EventsRepo
	UpdatedEventRepo     repository.EventsUpdateRepo
	UserRepo             repository.UsersRepo
	ProfileRepo          repository.OrganizerProfileRepo
	CategoryRepo         repository.CategoriesRepo
	EventExpiryPublisher NewTaskEventExpiryPublisher
	EmailTaskPublisher   NewTaskEmail
	rdb                  *redis.Client
}

func NewEventService(
	eventRepo repository.EventsRepo,
	updatedEventRepo repository.EventsUpdateRepo,
	userRepo repository.UsersRepo,
	profileRepo repository.OrganizerProfileRepo,
	categoryRepo repository.CategoriesRepo,
	eventExpiryPublisher NewTaskEventExpiryPublisher,
	emailTaskPublisher NewTaskEmail,
	rdb *redis.Client,
) *EventService {
	return &EventService{
		EventRepo:            eventRepo,
		UpdatedEventRepo:     updatedEventRepo,
		UserRepo:             userRepo,
		ProfileRepo:          profileRepo,
		CategoryRepo:         categoryRepo,
		EventExpiryPublisher: eventExpiryPublisher,
		EmailTaskPublisher:   emailTaskPublisher,
		rdb:                  rdb,
	}
}

var (
	eventBannerPath        = "./storage/event/banner"
	updatedEventBannerPath = "./storage/updated/banner"
)

var eventCache []string = []string{
	"events:all:*",
	"events:nearest:*",
	"organizer_events:all:*",
}

func (s *EventService) CreateEvent(banner *multipart.FileHeader, req *dto.EventReq) error {
	// Search the eo profile
	profile, err := s.ProfileRepo.FindByUserID(req.UserID)
	if err != nil {
		return errors.New("profile not found")
	}

	// Validate len categories
	if len(req.Categories) == 0 {
		return errors.New("categories or ticket cannot be empty")
	}

	// Get categories
	categories, err := s.CategoryRepo.FindByIDs(req.Categories)
	if err != nil {
		return err
	}

	// Convert unix to date
	startTime := helper.ConvertUnixtoDate(req.StartTime)
	endTime := helper.ConvertUnixtoDate(req.EndTime)

	// Default status
	status := string(model.Pending)
	if req.Status != "" {
		status = req.Status
	}

	// Get coordinates from req
	location := getLocation(req.Address.Lat, req.Address.Long)
	if location.Err != nil {
		return *location.Err
	}

	// Save banner to storage
	var eventBanner *string
	if banner != nil {
		_, fileName, err := helper.SaveToLocal(banner, eventBannerPath)
		if err != nil {
			log.Printf("error save image")
			return err
		}

		eventBanner = &fileName
	}

	// Save events
	event := &model.Events{
		ProfileID:     profile.ID,
		Banner:        eventBanner,
		Name:          req.Name,
		Slug:          utils.CreateSlug(req.Name),
		Status:        status,
		Description:   req.Description,
		Address:       *location.Address,
		City:          *location.City,
		Country:       *location.Country,
		DetailAddress: req.Address.DetailAddress,
		Coordinates:   *location.Coordinates,
		StartTime:     startTime.UTC(),
		EndTime:       endTime.UTC(),
	}

	event, err = s.EventRepo.Create(event, categories)
	if err != nil {
		log.Printf("error creating event")
		return err
	}

	// Invalidate cache after update
	utils.InvalidateCache(s.rdb, eventCache)

	// Email the admins
	// Only if the organizer decide to immediately up the event
	if event.Status == string(model.Pending) {
		admins, err := s.UserRepo.FindByRole(string(model.Admin))
		if err != nil {
			log.Println("[ERROR] admin data not found")
		}

		for _, admin := range admins {
			AdminEmailPayload := &model.EventEmailPayload{
				To:        admin.Email,
				EOName:    profile.Name,
				EOEmail:   profile.User.Email,
				EventName: event.Name,
				Status:    string(model.Create),
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
	}

	return nil
}

func (s *EventService) ReviewEvent(req *dto.ReviewEventReq) error {
	event, err := s.EventRepo.FindByID(req.ID)
	if err != nil {
		return errors.New("event not found")
	}

	if event.Status != string(model.Pending) && event.Status != string(model.Rejected) {
		return errors.New(fmt.Sprintf("event status is %s", event.Status))
	}

	now := time.Now().UTC()
	event.Status = req.Status
	event.ReviewedBy = req.ReviewedBy
	event.ReviewedAt = &now

	if req.Status == "reject" {
		event.RejectedReason = req.Reason
	} else {
		event.RejectedReason = nil
	}

	if err := s.EventRepo.ReviewEvent(event); err != nil {
		return err
	}

	if req.Status == string(model.Active) {
		payload := &model.EventExpiredPayload{EventID: event.ID}
		if err := s.EventExpiryPublisher.EnqueueEventExpiry(payload, event.EndTime); err != nil {
			log.Printf("[EXPIRY] failed to enqueue expiry for event %s: %v", event.ID, err)
		}
	}

	utils.InvalidateCache(s.rdb, eventCache)

	// Email the EO
	organizer, err := s.ProfileRepo.FindByID(event.ProfileID)
	if err != nil {
		log.Println("[ERROR] organizer data not found")
	}

	EOEmailPayload := &model.EventEmailPayload{
		To:        organizer.User.Email,
		EventName: event.Name,
		Status:    req.Status,
		Reason:    helper.StringValue(req.Reason),
	}
	if err := s.EmailTaskPublisher.Enqueue(model.TypeEventEOVerification, EOEmailPayload); err != nil {
		log.Printf("[EMAIL] failed sending email to event organizer %s\n", organizer.User.Email)
	}

	return nil
}

func (s *EventService) GetEvents(filter *dto.EventFilter, pagination model.Pagination) (*model.PaginationRow[*dto.EventsResp], error) {
	var events *model.PaginationRow[*dto.EventsResp]

	filterBytes, _ := json.Marshal(filter)
	hash := sha1.Sum(filterBytes)
	filterHash := hex.EncodeToString(hash[:])

	cacheKey := fmt.Sprintf("events:all:%d:%d:%s:%s", pagination.Page, pagination.Limit, pagination.Sort, filterHash)

	// Try to get from cache
	val, err := s.rdb.Get(context.Background(), cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &events)
	}

	if events == nil {
		// If cache miss, get from db
		events, err = s.EventRepo.FindAll(filter, pagination)
		if err != nil {
			return nil, err
		}

		// Set cache with 15 minute TTL
		if data, err := json.Marshal(events); err == nil {
			s.rdb.Set(context.Background(), cacheKey, data, 15*time.Minute)
		}
	}

	return events, nil
}

func (s *EventService) GetEventsByProfileID(userID string, filter *dto.EventFilter, pagination model.Pagination) (*model.PaginationRow[*dto.EventsResp], error) {
	profile, err := s.ProfileRepo.FindByUserID(userID)
	if err != nil {
		return nil, errors.New("profile not found")
	}

	filter.ProfileID = &profile.ID

	var events *model.PaginationRow[*dto.EventsResp]

	filterBytes, _ := json.Marshal(filter)
	hash := sha1.Sum(filterBytes)
	filterHash := hex.EncodeToString(hash[:])

	cacheKey := fmt.Sprintf("organizer_events:all:%d:%d:%s:%s:%s", pagination.Page, pagination.Limit, pagination.Sort, filterHash)

	// Try to get from cache
	val, err := s.rdb.Get(context.Background(), cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &events)
	}

	if events == nil {
		// If cache miss, get from db
		events, err = s.EventRepo.FindByProfileID(filter, pagination)
		if err != nil {
			return nil, err
		}

		// Set cache with 15 minute TTL
		if data, err := json.Marshal(events); err == nil {
			s.rdb.Set(context.Background(), cacheKey, data, 15*time.Minute)
		}
	}

	return events, nil
}

func (s *EventService) GetEventByID(id string, userLat, userLon float64) (*dto.EventsResp, error) {
	// Search the event
	event, err := s.EventRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("event not found")
	}

	event.Banner = helper.StrPointerIfNotEmpty(
		func() string {
			if event.Banner == nil {
				return ""
			}
			return fmt.Sprintf("http://localhost:8080/api/v1/event/banner/%s", *event.Banner)
		}(),
	)

	// Get the organizer profile
	organizer, err := s.ProfileRepo.FindByID(event.ProfileID)
	if err != nil {
		return nil, errors.New("organizer not found")
	}

	user, err := s.UserRepo.FindByID(organizer.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Only profile that have the event can access on the organizer dashboard
	// and also admin can access the event, check if the organizer profile id is equal to the event profile id
	if organizer.ID != event.ProfileID && helper.StringValue(user.Role) != string(model.Admin) {
		return nil, errors.New("unauthorized access to event")
	}

	organizer.PhotoProfile = helper.StrPointerIfNotEmpty(
		func() string {
			if organizer.PhotoProfile == nil {
				return ""
			}
			return fmt.Sprintf("http://localhost:8080/api/v1/profile/photo/%s", *organizer.PhotoProfile)
		}(),
	)

	var eventCategories []dto.EventCategories
	for _, category := range event.Categories {
		eventCategories = append(eventCategories, dto.EventCategories{
			ID:   category.Category.ID,
			Name: category.Category.Name,
		})
	}

	havDist := utils.Haversine(userLat, userLon, event.Lat, event.Lon)
	distance := fmt.Sprintf("%.2f km", havDist)

	respReq := &dto.EventRespReq{
		Event:           event,
		Organizer:       organizer,
		EventCategories: eventCategories,
		StartTime:       helper.ConvertDatetoUnix(event.StartTime.Format(time.RFC3339)),
		EndTime:         helper.ConvertDatetoUnix(event.EndTime.Format(time.RFC3339)),
		Distance:        distance,
		Path: []dto.PathPoint{
			{Name: "user", Lat: userLat, Lon: userLon},
			{Name: event.Name, Lat: event.Lat, Lon: event.Lon},
		},
		CreatedAt: helper.ConvertDatetoUnix(event.CreatedAt.Format(time.RFC3339)),
		UpdatedAt: helper.ConvertDatetoUnix(event.UpdatedAt.Format(time.RFC3339)),
		DeletedAt: helper.TimePtrToUnix(event.DeletedAt),
	}
	eventResp, err := dto.ToEventResp(respReq)
	if err != nil {
		return nil, err
	}

	return eventResp, nil
}

func (s *EventService) GetEventRoute(id string, userLat, userLon float64) (*dto.RouteResp, error) {
	event, err := s.EventRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("event not found")
	}

	distance, path := utils.ComputePathToEvent(userLat, userLon, event.Name, event.Lat, event.Lon)

	// Compute analytic in background
	// It will display the comparison between Dijkstra and Haversine on calculate the distance
	go func() {
		user := model.Location{Name: "user", Lat: userLat, Lon: userLon}
		eventLoc := model.Location{Name: event.Name, Lat: event.Lat, Lon: event.Lon}
		events := []model.Location{eventLoc}

		go analytics.ComputeAnalytic(user, events)
	}()

	return &dto.RouteResp{
		Event:    event.Name,
		Distance: distance,
		Path:     path,
	}, nil
}

func (s *EventService) FindNearestEvent(user model.Location, pagination model.Pagination) (*model.PaginationRow[*dto.EventsResp], error) {
	var events *model.PaginationRow[*dto.EventsResp]

	cacheKey := fmt.Sprintf("events:nearest:%d:%d:%d:%f:%f", pagination.Page, pagination.Limit, pagination.Sort, user.Lat, user.Lon)

	// Try to get from cache
	val, err := s.rdb.Get(context.Background(), cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &events)
	}

	if events == nil {
		// If cache miss, get from db
		events, err = s.EventRepo.FindNearestEvents(user.Lat, user.Lon, pagination)
		if err != nil {
			return nil, errors.New("there are no events near your location at the moment.")
		}

		// Set cache with 15 minute TTL
		if data, err := json.Marshal(events); err == nil {
			s.rdb.Set(context.Background(), cacheKey, data, 15*time.Minute)
		}
	}

	return events, nil
}

func (s *EventService) UpdateEvent(banner *multipart.FileHeader, req *dto.EventReq) error {
	// Search the eo profile
	profile, err := s.ProfileRepo.FindByUserID(req.UserID)
	if err != nil {
		return errors.New("profile not found")
	}

	// Validate event
	event, err := s.EventRepo.FindByID(*req.ID)
	if err != nil {
		return errors.New("event not found")
	}

	originalStatus := event.Status

	// Validate the user
	if !helper.IsAuthorized(event.ProfileID, profile.ID) {
		return errors.New("unauthorized action")
	}

	// Event with status pending cannot be updated
	if event.Status == string(model.Pending) {
		return errors.New(fmt.Sprintf("event status is %s", event.Status))
	}

	// If event status already active
	// The organizer can't update it to draft
	if event.Status == string(model.Active) && req.Status == string(model.Draft) {
		return errors.New(fmt.Sprintf("An active event cannot be reverted to %s status", req.Status))
	}

	// Check if its critical changed
	// It only happen if the event already approve by the admin
	if event.Status == string(model.Active) {
		if err := s.CreateUpdateEvent(banner, event, req); err != nil {
			return err
		}

		// Invalidate cache after update
		utils.InvalidateCache(s.rdb, updatedEventCache)

		// Email the admins
		admins, err := s.UserRepo.FindByRole(string(model.Admin))
		if err != nil {
			log.Println("[ERROR] admin data not found")
		}

		for _, admin := range admins {
			AdminEmailPayload := &model.EventEmailPayload{
				To:        admin.Email,
				EOName:    profile.Name,
				EOEmail:   profile.User.Email,
				EventName: event.Name,
				Status:    string(model.Update),
			}
			if err := s.EmailTaskPublisher.Enqueue(model.TypeEventAdminNotification, AdminEmailPayload); err != nil {
				log.Printf("[EMAIL] failed sending email to admin %s\n", admin.Email)
			}
		}

		return nil
	}

	// Get categories
	categories, err := s.CategoryRepo.FindByIDs(req.Categories)
	if err != nil {
		return err
	}

	// Convert unix to date
	startTime := helper.ConvertUnixtoDate(req.StartTime)
	endTime := helper.ConvertUnixtoDate(req.EndTime)

	// Get coordinates from req
	location := getLocation(req.Address.Lat, req.Address.Long)
	if location.Err != nil {
		return *location.Err
	}

	// Save events
	if req.Status != "" {
		event.Status = req.Status
	}

	// if the event owner decide to publish draft event
	// Check the banner is uploaded or not
	if event.Status == string(model.Pending) && banner == nil {
		return errors.New("banner is required")
	}

	// If banner changed
	// Save  to storage
	oldBanner := event.Banner
	if banner != nil {
		_, fileName, err := helper.SaveToLocal(banner, eventBannerPath)
		if err != nil {
			log.Printf("[ERROR] failed to save banner %v\n", err)
			return err
		}
		event.Banner = &fileName

		// Delete old banner from storage
		if oldBanner != nil {
			os.Remove(*oldBanner)
		}
	}

	event.Name = req.Name
	event.Slug = utils.CreateSlug(req.Name)
	event.Description = req.Description
	event.Address = *location.Address
	event.City = *location.City
	event.Country = *location.Country
	event.DetailAddress = req.Address.DetailAddress
	event.Coordinates = *location.Coordinates
	event.StartTime = startTime
	event.EndTime = endTime

	if err := s.EventRepo.Update(event, categories); err != nil {
		return err
	}

	// Invalidate cache after update
	utils.InvalidateCache(s.rdb, eventCache)

	// If organizer decide to up the event
	// Notify the admins
	if event.Status == string(model.Pending) {
		admins, err := s.UserRepo.FindByRole(string(model.Admin))
		if err != nil {
			log.Println("[ERROR] admin data not found")
		}

		statusStr := string(model.Update)
		if originalStatus == string(model.Draft) {
			statusStr = string(model.Create)
		}

		for _, admin := range admins {
			AdminEmailPayload := &model.EventEmailPayload{
				To:        admin.Email,
				EOName:    profile.Name,
				EOEmail:   profile.User.Email,
				EventName: event.Name,
				Status:    statusStr,
			}
			if err := s.EmailTaskPublisher.Enqueue(model.TypeEventAdminNotification, AdminEmailPayload); err != nil {
				log.Printf("[EMAIL] failed sending email to admin %s\n", admin.Email)
			}
		}
	}

	return nil
}

func (s *EventService) CancelEvent(id, userID string) error {
	// Find organizer
	organizer, err := s.ProfileRepo.FindByUserID(userID)
	if err != nil {
		return errors.New("organizer not found")
	}

	// Validate event
	event, err := s.EventRepo.FindByID(id)
	if err != nil {
		return errors.New("event not found")
	}

	// Validate organizer
	if !helper.IsAuthorized(event.ProfileID, organizer.ID) {
		return errors.New("unauthorized action")
	}

	// Only event with status active can be canceled
	if event.Status != string(model.Active) {
		return errors.New(fmt.Sprintf("event status is %s", event.Status))
	}

	// Validate cancelation
	// Event can be canceled less then 3 days before it starts
	if time.Until(event.StartTime) < 72*time.Hour {
		return errors.New("event cannot be cancelled less than 3 days before it starts")
	}

	// Cancel event
	if err := s.EventRepo.CancelEvent(event.ID); err != nil {
		log.Printf("[ERROR] error canceling event %v\n", err)
		return errors.New("failed to cancel")
	}

	if err := s.EventExpiryPublisher.CancelEventExpiry(event.ID); err != nil {
		log.Printf("[EXPIRY] failed to cancel expiry task for event %s: %v", event.ID, err)
	}

	// Invalidate cache after update
	utils.InvalidateCache(s.rdb, eventCache)

	return nil
}

func (s *EventService) DeleteEvent(id, userID string) error {
	event, err := s.EventRepo.FindByID(id)
	if err != nil {
		return errors.New("Event not found")
	}

	if event.Profile.User.ID != userID {
		return errors.New("Unauthorized action")
	}

	// Only delete event with status draft
	if event.Status == string(model.Draft) {
		return s.EventRepo.Delete(id)
	}

	utils.InvalidateCache(s.rdb, eventCache)

	return nil
}

func (s *EventService) CreateUpdateEvent(banner *multipart.FileHeader, event *model.Events, req *dto.EventReq) error {
	// Update only permitted 1 week before the event
	if time.Until(event.StartTime) < 7*24*time.Hour {
		return errors.New("event cannot be updated within 7 days of the event date")
	}

	// Validate len categories & tickets
	if len(req.Categories) == 0 {
		return errors.New("categories cannot be empty")
	}

	categories, err := s.CategoryRepo.FindByIDs(req.Categories)
	if err != nil {
		return errors.New("category not found")
	}

	// Get coordinates from req
	location := getLocation(req.Address.Lat, req.Address.Long)
	if location.Err != nil {
		return *location.Err
	}

	// If banner changed
	// Save to temporary storage
	var updatedEvent *model.UpdatedEvents
	if banner != nil {
		_, fileName, err := helper.SaveToLocal(banner, updatedEventBannerPath)
		if err != nil {
			log.Printf("[ERROR] failed to save banner %v\n", err)
			return err
		}

		updatedEvent = &model.UpdatedEvents{
			EventID:       event.ID,
			Name:          req.Name,
			Banner:        &fileName,
			Slug:          utils.CreateSlug(req.Name),
			Status:        string(model.Pending),
			Description:   req.Description,
			Address:       *location.Address,
			City:          *location.City,
			Country:       *location.Country,
			DetailAddress: req.Address.DetailAddress,
			Coordinates:   *location.Coordinates,
			StartTime:     helper.ConvertUnixtoDate(req.StartTime),
			EndTime:       helper.ConvertUnixtoDate(req.EndTime),
		}

		if err := s.UpdatedEventRepo.Create(updatedEvent, categories); err != nil {
			log.Printf("[ERROR] failed to update event %v\n", err)
			return err
		}
	} else {
		fileName := *event.Banner

		eventBannerSrc := filepath.Join(eventBannerPath, fileName)
		dstPath := filepath.Join(eventBannerPath, fileName)

		bannerFile, err := helper.CopyFile(eventBannerSrc, dstPath)
		if err != nil {
			return err
		}

		updatedEvent = &model.UpdatedEvents{
			EventID:       event.ID,
			Name:          req.Name,
			Banner:        &bannerFile,
			Slug:          utils.CreateSlug(req.Name),
			Status:        string(model.Pending),
			Description:   req.Description,
			Address:       *location.Address,
			City:          *location.City,
			Country:       *location.Country,
			DetailAddress: req.Address.DetailAddress,
			Coordinates:   *location.Coordinates,
			StartTime:     helper.ConvertUnixtoDate(req.StartTime),
			EndTime:       helper.ConvertUnixtoDate(req.EndTime),
		}

		if err := s.UpdatedEventRepo.Create(updatedEvent, categories); err != nil {
			log.Printf("[ERROR] failed to update event %v\n", err)
			return err
		}
	}

	return nil
}

func getLocation(lat, lon string) *dto.LocationResp {
	coordinates := fmt.Sprintf("POINT(%s %s)", lon, lat)
	reverseResp, err := utils.ReverseGeocode(lat, lon)
	if err != nil {
		log.Printf(err.Error())
		return &dto.LocationResp{
			Coordinates: nil,
			Address:     nil,
			City:        nil,
			Country:     nil,
			Err:         &err,
		}
	}

	address := reverseResp.DisplayName

	city := reverseResp.Address.City

	if city == "" {
		city = reverseResp.Address.Town
	}

	if city == "" {
		city = reverseResp.Address.Village
	}

	country := reverseResp.Address.Country

	return &dto.LocationResp{
		Coordinates: &coordinates,
		Address:     &address,
		City:        &city,
		Country:     &country,
		Err:         nil,
	}
}
