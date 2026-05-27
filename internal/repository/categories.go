package repository

import (
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"strings"

	"gorm.io/gorm"
)

type CategoriesRepository struct {
	db *gorm.DB
}

// GetWithPagination implements [CategoriesRepo].
func (r *CategoriesRepository) GetWithPagination(
	pagination model.Pagination,
	filter *dto.FilterCatReq,
) (*model.PaginationRow[*dto.ListCatResp], error) {

	var categories []*dto.ListCatResp

	query := r.db.
		Table("categories").
		Select(`
			categories.id,
			categories.name,
			categories.slug,
			categories.created_at,
			categories.updated_at,
			COUNT(event_categories.category_id) AS total_used
		`).
		Joins(`
			LEFT JOIN event_categories
			ON categories.id = event_categories.category_id
			AND event_categories.deleted_at IS NULL
		`).
		Scopes(FilterCat(filter)).
		Group("categories.id")

	if err := query.
		Scopes(Paginate(categories, &pagination, query)).
		Order("total_used DESC").
		Scan(&categories).Error; err != nil {
		return nil, err
	}

	return &model.PaginationRow[*dto.ListCatResp]{
		Pagination: pagination,
		Rows:       categories,
	}, nil
}

// Create implements CategoriesRepo.
func (r *CategoriesRepository) Create(category *model.Categories) error {
	return r.db.Create(category).Error
}

// Delete implements CategoriesRepo.
func (r *CategoriesRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.Categories{}).Error
}

// FindAll implements CategoriesRepo.
func (r *CategoriesRepository) FindAll() ([]*model.Categories, error) {
	var categories []*model.Categories

	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

// FindByID implements CategoriesRepo.
func (r *CategoriesRepository) FindByID(id string) (*model.Categories, error) {
	var category *model.Categories

	if err := r.db.Where("id = ?", id).First(&category).Error; err != nil {
		return nil, err
	}

	return category, nil
}

// FindByIDs implements CategoriesRepo.
func (r *CategoriesRepository) FindByIDs(ids []int64) ([]*model.Categories, error) {
	var categories []*model.Categories

	if err := r.db.Where("(id) IN ?", ids).Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

// FindBySlug implements CategoriesRepo.
func (r *CategoriesRepository) FindBySlug(title string, pagination model.Pagination) (*model.PaginationRow[*model.Categories], error) {
	var categories []*model.Categories

	query := r.db.Where("LOWER(slug) LIKE LOWER(?)", "%"+title+"%")

	if err := query.Scopes(Paginate(categories, &pagination, query)).Find(&categories).Error; err != nil {
		return nil, err
	}

	return &model.PaginationRow[*model.Categories]{
		Pagination: pagination,
		Rows:       categories,
	}, nil
}

func FilterCat(filter *dto.FilterCatReq) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if filter.Name != nil {
			query := "%" + strings.ToLower(*filter.Name) + "%"
			db = db.Where("LOWER(slug) LIKE LOWER(?)", query)
		}
		return db
	}
}

// Update implements CategoriesRepo.
func (r *CategoriesRepository) Update(updateCat *model.Categories) error {
	return r.db.Save(updateCat).Error
}

func NewCategoriesRepository(db *gorm.DB) CategoriesRepo {
	return &CategoriesRepository{db: db}
}
