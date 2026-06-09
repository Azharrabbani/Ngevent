package dto

type RegisterInput struct {
	Email           string `json:"email" validate:"required"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

type RoleInput struct {
	Role string `json:"role" validate:"required,oneof=user 'event organizer'"`
}

type LoginInput struct {
	Email      string `json:"email" validate:"required"`
	Password   string `json:"password" validate:"required"`
	RememberMe bool   `json:"remember_me"`
}

type ForgetPasswordInput struct {
	Email string `json:"email" validate:"required"`
}

type ResetPasswordInput struct {
	NewPassword     string `json:"new_password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

type VerifyEmailInput struct {
	Email string `json:"email"`
	OTP   string `json:"otp" validate:"required"`
}

type ResentOTPInput struct {
	Email string `json:"email" validate:"required"`
}

type ListUsersReq struct {
	Role       *string `json:"role" query:"role"`
	IsVerified *bool   `json:"is_verified" query:"is_verified"`
	Email      *string `json:"email" query:"email"`
}

type LoginResponse struct {
	ID              string  `json:"id"`
	Email           string  `json:"email"`
	Role            *string `json:"role"`
	NgeventToken    string  `json:"ngevent-token"`
	NgeventRefToken string  `json:"ngevent-ref-token"`
}

type UsersResponse struct {
	ID         string  `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email      string  `json:"email"`
	Role       *string `json:"role"`
	HasProfile *bool   `json:"has_profile,omitempty"`
	IsVerified bool    `json:"is_verified"`
	CreatedAt  int64   `json:"created_at"`
	UpdatedAt  int64   `json:"updated_at"`
	DeletedAt  int64   `json:"deleted_at,omitempty"`
}

type RefreshTokenResp struct {
	NgeventToken string `json:"ngevent-token"`
}
