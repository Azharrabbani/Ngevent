package service

import (
	"errors"
	"ngevent/internal/dto"
	"ngevent/internal/repository"
	"ngevent/internal/utils/helper"

	"github.com/redis/go-redis/v9"
)

type EventService struct {
	EventRepo repository.EventsRepo
	rdb       *redis.Client
}

func NewEventService(
	eventRepo repository.EventsRepo,
	rdb *redis.Client,
) *EventService {
	return &EventService{
		EventRepo: eventRepo,
	}
}

func (s *EventService) CreateEvent(req *dto.EventReq) error {
	// Validate len categories & tickets
	if len(req.Categories) < 0 || len(req.Tickets) < 0 {
		return errors.New("categories or ticket cannot be empty")
	}

	// Convert unix to date
	date := helper.ConvertUnixtoDate(req.Date)
	
	// Save events

}
