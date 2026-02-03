package dto

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
