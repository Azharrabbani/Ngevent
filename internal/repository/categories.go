package repository

import (
	"ngevent/internal/model"

	"gorm.io/gorm"
)

type CategoriesRepository struct {
	db *gorm.DB
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
func (r *CategoriesRepository) FindAll(pagination model.Pagination) (*model.PaginationRow[*model.Categories], error) {
	var categories []*model.Categories

	if err := r.db.Scopes(Paginate(categories, &pagination, r.db)).Find(&categories).Error; err != nil {
		return nil, err
	}

	return &model.PaginationRow[*model.Categories]{
		Pagination: pagination,
		Rows:       categories,
	}, nil
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

// Update implements CategoriesRepo.
func (r *CategoriesRepository) Update(updateCat *model.Categories) error {
	return r.db.Save(updateCat).Error
}

func NewCategoriesRepository(db *gorm.DB) CategoriesRepo {
	return &CategoriesRepository{db: db}
}
