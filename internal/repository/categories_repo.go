package repository

import (
	"ngevent/internal/dto"
	"ngevent/internal/model"
)

type CategoriesRepo interface {
	Create(category *model.Categories) error
	FindAll() ([]*model.Categories, error)
	GetWithPagination(pagination model.Pagination, filter *dto.FilterCatReq) (*model.PaginationRow[*dto.ListCatResp], error)
	FindByID(id string) (*model.Categories, error)
	FindByIDs(ids []int64) ([]*model.Categories, error)
	FindBySlug(title string, pagination model.Pagination) (*model.PaginationRow[*model.Categories], error)
	Update(updateCat *model.Categories) error
	Delete(id string) error
}
