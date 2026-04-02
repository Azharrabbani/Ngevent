package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"ngevent/internal/repository"
	"ngevent/internal/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

type CategoryService struct {
	CategoryRepo repository.CategoriesRepo
	rdb          *redis.Client
}

func NewCategoryService(
	categoryRepo repository.CategoriesRepo,
	rdb *redis.Client,
) *CategoryService {
	return &CategoryService{
		CategoryRepo: categoryRepo,
		rdb:          rdb,
	}
}

var categoryCache []string = []string{
	"category:all:*",
}

func (s *CategoryService) Create(req *dto.CreateCatReq) error {
	category := &model.Categories{
		Name: req.Name,
		Slug: utils.CreateSlug(req.Name),
	}

	if err := s.CategoryRepo.Create(category); err != nil {
		return err
	}

	// Invalidate cache after update
	utils.InvalidateCache(s.rdb, categoryCache)

	return nil
}

func (s *CategoryService) FindAll(pagination model.Pagination) (*model.PaginationRow[*model.Categories], error) {
	var categories *model.PaginationRow[*model.Categories]

	// Generate cache key
	cachekey := fmt.Sprintf("category:all:%d:%d:%s", pagination.Page, pagination.Limit, pagination.Sort)

	// Try to get from cache
	val, err := s.rdb.Get(context.Background(), cachekey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &categories)
	}

	if categories == nil {
		// If cache miss, get from db
		categories, err = s.CategoryRepo.FindAll(pagination)
		if err != nil {
			return nil, err
		}

		// Set cache with 15 minute TTL
		if data, err := json.Marshal(categories); err == nil {
			s.rdb.Set(context.Background(), cachekey, data, 15*time.Minute)
		}
	}

	return categories, nil
}

func (s *CategoryService) FindByID(id string) (*model.Categories, error) {
	return s.CategoryRepo.FindByID(id)
}

func (s *CategoryService) FindBySlug(title string, pagination model.Pagination) (*model.PaginationRow[*model.Categories], error) {
	var categories *model.PaginationRow[*model.Categories]
	title = utils.CreateSlug(title)

	// Generate cache key
	cacheKey := fmt.Sprintf("category:all:%s:%d:%d:%s", title, pagination.Page, pagination.Limit, pagination.Sort)

	// Try get from cache
	val, err := s.rdb.Get(context.Background(), cacheKey).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &categories)
	}

	if categories == nil {
		// If miss get from db
		categories, err = s.CategoryRepo.FindBySlug(title, pagination)
		if err != nil {
			return nil, err
		}

		// Set cache with 15 minute TTL
		if data, err := json.Marshal(categories); err == nil {
			s.rdb.Set(context.Background(), cacheKey, data, 15*time.Minute)
		}

	}

	return categories, nil
}

func (s *CategoryService) Update(updateReq *dto.UpdateCatReq) error {
	// Search the category
	category, err := s.CategoryRepo.FindByID(updateReq.CategoryID)
	if err != nil {
		return errors.New("category not found")
	}

	// Update the category
	category.Name = updateReq.Name
	category.Slug = utils.CreateSlug(updateReq.Name)

	// Save the update
	if err := s.CategoryRepo.Update(category); err != nil {
		return err
	}

	// Invalidate cache after update
	utils.InvalidateCache(s.rdb, categoryCache)

	return nil
}

func (s *CategoryService) Delete(id string) error {
	// Search the category
	_, err := s.CategoryRepo.FindByID(id)
	if err != nil {
		return errors.New("category not found")
	}

	// Delete the category
	if err := s.CategoryRepo.Delete(id); err != nil {
		return err
	}

	// Invalidate cache after update
	utils.InvalidateCache(s.rdb, categoryCache)

	return nil
}
