package dto

type PaginationResponse[T any] struct {
	Data       []T `json:"data"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalData  int `json:"total_data"`
	TotalPages int `json:"total_pages"`
}
