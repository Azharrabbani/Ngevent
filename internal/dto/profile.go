package dto

import "mime/multipart"

// =========== Attendee profile req ================
type CreateAttendeeProfileReq struct {
	UserID       string
	PhotoProfile *multipart.FileHeader
	Name         string
	Username     *string
	PhoneNumber  string
	ISO          string
	Address      *string
}

type UpdateAttendeeProfileReq struct {
	Name        string  `json:"name" validate:"required"`
	Username    *string `json:"username"`
	PhoneNumber string  `json:"phone_number" validate:"required"`
	ISO         string  `json:"iso" validate:"required"`
	Address     *string `json:"address"`
}

// =========== Organizer Profile req ============
type CreateOrganizerProfileReq struct {
	UserID        string
	PhotoProfile  *multipart.FileHeader
	Name          string
	PhoneNumber   string
	ISO           string
	Address       *string
	SocialMedia   OrganizerSocialMediaReq
	CompanyDetail OrganizerCompDetailReq
}

type OrganizerSocialMediaReq struct {
	Email     *string `json:"email,omitempty"`
	Instagram *string `json:"instagram,omitempty"`
}

type OrganizerCompDetailReq struct {
	Description *string `json:"description,omitempty"`
	NPWP        string  `json:"npwp" validate:"required"`
	NIB         string  `json:"nib" validate:"required"`
}

type UpdateOrganizerProfileReq struct {
	Name          string                  `json:"name" validate:"required"`
	PhoneNumber   string                  `json:"phone_number" validate:"required"`
	ISO           string                  `json:"iso" validate:"required"`
	Address       *string                 `json:"address"`
	SocialMedia   OrganizerSocialMediaReq `json:"social_media"`
	CompanyDetail OrganizerCompDetailReq  `json:"company_detail"`
}

// =========== Attendee profile response ==============
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

// =========== Organizer profile response ==============
type OrganizerProfilesResponse struct {
	ID            string                  `json:"id"`
	UserID        string                  `json:"user_id"`
	IsVerified    bool                    `json:"is_verified"`
	Email         string                  `json:"email"`
	Name          string                  `json:"name"`
	PhotoProfile  *string                 `json:"photo_profile,omitempty"`
	PhoneNumber   string                  `json:"phone_number"`
	Country       string                  `json:"country"`
	Address       *string                 `json:"address,omitempty"`
	SocialMedia   OrganizerSocialMediaReq `json:"social_media"`
	CompanyDetail OrganizerCompDetailReq  `json:"company_detail"`
}
