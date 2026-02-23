package dto

type RegisterInput struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=admin user 'event organizer'"`
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
	OTP string `json:"otp" validate:"required"`
}

type ResentOTPInput struct {
	Email string `json:"email" validate:"required"`
}

type LoginResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Role        string `json:"role"`
}
