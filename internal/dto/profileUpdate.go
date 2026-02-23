package dto

type ValidateUpdateReq struct {
	Status   string `json:"status" validate:"required,oneof=pending approved rejected"`
	UpdateID string
	Reason   string `json:"reason"`
}
