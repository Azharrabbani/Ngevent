package dto

type RegisterInput struct {
	Email string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Role string `json:"role" validate:"required,oneof=admin user 'event organizer'"`
}

type LoginInput struct {
	Email string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}