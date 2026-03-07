package repository

import "ngevent/internal/model"

type CategoriesRepo interface {
	Create(category *model.Categories) error
	FindAll(pagination model.Pagination) (*model.PaginationRow[*model.Categories], error)
	FindByID(id string) (*model.Categories, error)
	FindBySlug(title string, pagination model.Pagination) (*model.PaginationRow[*model.Categories], error)
	Update(updateCat *model.Categories) error
	Delete(id string) error
}
