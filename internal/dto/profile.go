package dto

type CreateProfileReq struct {
	UserID      string
	Name        string
	Username    *string
	PhoneNumber string
	ISO         string
	Address     *string
}

type UpdateProfileReq struct {
	UserID      string  `json:"user_id"`
	Name        string  `json:"name"`
	Username    *string `json:"username"`
	PhoneNumber string  `json:"phone_number"`
	ISO         string  `json:"iso"`
	Address     *string `json:"address"`
}

type AttendeeProfilesResponse struct {
	ID           string  `json:"id"`
	UserID       string  `json:"user_id"`
	Email        string  `json:"email"`
	Name         string  `json:"name"`
	Username     *string `json:"username,omitempty"`
	PhotoProfile *string `json:"photo_profile,omitempty"`
	PhoneNumber  string  `json:"phone_number"`
	Country      string  `json:"country"`
	Address      *string `json:"address,omitempty"`
}
