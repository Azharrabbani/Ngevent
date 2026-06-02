package dto

import (
	"mime/multipart"
	"time"
)

// =========== Attendee profile req ================
type CreateAttendeeProfileReq struct {
	UserID       string
	PhotoProfile *multipart.FileHeader
	Name         string  `json:"name" validate:"required"`
	Username     *string `json:"username"`
	PhoneNumber  string  `json:"phone_number" validate:"required"`
	ISO          string
	Address      *string `json:"address"`
}

type FilterProfileReq struct {
	Filter *string `json:"filter" query:"filter"`
	Status *string `json:"status" query:"status"`
}

type FilterPublicProfileReq struct {
	Filter *string `json:"filter" query:"filter"`
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
	Name          string `form:"name" validate:"required"`
	PhoneNumber   string `form:"phone_number" validate:"required"`
	ISO           string
	Address       *string                 `form:"address" validate:"required"`
	SocialMedia   OrganizerSocialMediaReq `validate:"required"`
	CompanyDetail OrganizerCompDetailReq  `validate:"required"`
}

type OrganizerSocialMediaReq struct {
	Email     *string `json:"email,omitempty" form:"email,omitempty"`
	Instagram *string `json:"instagram,omitempty" form:"instagram,omitempty"`
}

type FilterReq struct {
	Country string `json:"country" query:"country" validate:"required"`
}

type OrganizerCompDetailReq struct {
	Description *string               `form:"description,omitempty"`
	NPWP        string                `form:"npwp" validate:"required"`
	NPWPFile    *multipart.FileHeader `form:"npwp_file"`
	NIB         string                `form:"nib" validate:"required"`
	NIBFile     *multipart.FileHeader `form:"nib_file"`
}

type UpdateOrganizerProfileReq struct {
	Name          string                  `json:"name" validate:"required"`
	PhoneNumber   string                  `json:"phone_number" validate:"required"`
	ISO           string                  `json:"iso" validate:"required"`
	Address       *string                 `json:"address"`
	SocialMedia   OrganizerSocialMediaReq `json:"social_media"`
	CompanyDetail OrganizerCompDetailReq  `json:"company_detail"`
}

type SaveNPWPAndNIBFileReq struct {
	NPWP     *multipart.FileHeader
	NIB      *multipart.FileHeader
	NPWPPath string
	NIBPath  string
}

type ValidateFileReq struct {
	Photo *multipart.FileHeader
	NPWP  multipart.FileHeader
	NIB   multipart.FileHeader
}

type ApprovedReq struct {
	ReviewedBy string
	ReviewedAt time.Time
}

type RejectedReq struct {
	Reason     string `json:"reason" validate:"required"`
	ReviewedBy string
	ReviewedAt time.Time
}

type OrganizerCompDetailRes struct {
	Description *string `json:"description,omitempty"`
	NPWP        string  `json:"npwp"`
	NPWPFile    string  `json:"npwp_file"`
	NIB         string  `json:"nib"`
	NIBFile     string  `json:"nib_file"`
}

type OrganizerStatusResp struct {
	Status         string  `json:"status"`
	RejectedReason *string `json:"rejected_reason,omitempty"`
	ReviewedBy     *string `json:"reviewed_by,omitempty"`
	ReviewedAt     *int64  `json:"reviewed_at,omitempty"`
}

// =========== Attendee profile response ==============
type AttendeeProfilesResponse struct {
	ID           string  `json:"id"`
	UserID       string  `json:"user_id"`
	Email        string  `json:"email"`
	Name         string  `json:"name"`
	Username     *string `json:"username,omitempty"`
	PhotoProfile string  `json:"photo_profile"`
	PhoneNumber  string  `json:"phone_number"`
	Country      string  `json:"country"`
	Address      *string `json:"address,omitempty"`
}

// =========== Organizer profile response ==============
type OrganizerProfilesResponse struct {
	ID            string                  `json:"id"`
	UserID        string                  `json:"user_id"`
	Status        OrganizerStatusResp     `json:"status,omitempty"`
	Email         string                  `json:"email"`
	Name          string                  `json:"name"`
	PhotoProfile  string                  `json:"photo_profile,omitempty"`
	PhoneNumber   string                  `json:"phone_number"`
	Country       string                  `json:"country"`
	Address       *string                 `json:"address,omitempty"`
	SocialMedia   OrganizerSocialMediaReq `json:"social_media"`
	CompanyDetail OrganizerCompDetailRes  `json:"company_detail"`
	CreatedAt     int64                   `json:"created_at"`
	UpdatedAt     int64                   `json:"updated_at"`
}

type OrganizerUpdatesResponse struct {
	ID           string `json:"id"`
	ProfileID    string `json:"profile_id"`
	Status       string `json:"status"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Description  string `json:"description"`
	Address      string `json:"address"`
	Instagram    string `json:"instagram"`
	PhoneNumber  string `json:"phone_number"`
	Country      string `json:"country"`
	NPWPNumber   string `json:"npwp_number"`
	NPWPDocument string `json:"npwp_document"`
	NIBNumber    string `json:"nib_number"`
	NIBDocument  string `json:"nib_document"`
	CreatedAt    int64  `json:"created_at"`
	UpdatedAt    int64  `json:"updated_at"`
}
