package service

import (
	"context"
	"encoding/json"
	"fmt"
	"ngevent/internal/dto"
	"ngevent/internal/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

type LocationService struct {
	rdb *redis.Client
}

var locationCache []string = []string{
	"location:all:*",
}

func NewLocationService(rdb *redis.Client) *LocationService {
	return &LocationService{rdb: rdb}
}

func (s *LocationService) SearchLocation(query string) (*[]dto.SearchResponse, error) {
	var locations *[]dto.SearchResponse

	cacheKey := fmt.Sprintf("location:all:%s", query)

	// Try to get from cache
	val, err := s.rdb.Get(context.Background(), cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &locations)
	}

	if locations == nil {
		// If cache miss, get from db
		locations, err = utils.SearchLocation(query)
		if err != nil {
			return nil, err
		}

		// Set cache with 15 minute TTL
		if data, err := json.Marshal(locations); err == nil {
			s.rdb.Set(context.Background(), cacheKey, data, 15*time.Minute)
		}
	}

	return locations, nil
}
