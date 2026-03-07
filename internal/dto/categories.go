package dto

type CreateCatReq struct {
	Name string `json:"name" validate:"required"`
}

type UpdateCatReq struct {
	Name       string `json:"name" validate:"required"`
	CategoryID string
}

type FindCatReq struct {
	Name string `json:"name" query:"name"`
}
