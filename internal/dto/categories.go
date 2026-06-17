package dto

import (
	"time"
)

type CreateCatReq struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
}

type UpdateCatReq struct {
	Name       string `json:"name" validate:"required,min=2,max=50"`
	CategoryID string
}

type FindCatReq struct {
	Name string `json:"name" query:"name"`
}

type FilterCatReq struct {
	Name *string `json:"name" query:"name"`
}

type ListCatResp struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	TotalUsed int64     `json:"total_used"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
