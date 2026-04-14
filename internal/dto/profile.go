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

type FilterReq struct {
	Country string `json:"country" query:"country" validate:"required"`
}

type OrganizerCompDetailReq struct {
	Description *string              `json:"description,omitempty"`
	NPWP        string               `json:"npwp" validate:"required"`
	NPWPFile    multipart.FileHeader `json:"npwp_file" validate:"required"`
	NIB         string               `json:"nib" validate:"required"`
	NIBFile     multipart.FileHeader `json:"nib_file" validate:"required"`
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
	PhotoProfile *string `json:"photo_profile,omitempty"`
	PhoneNumber  string  `json:"phone_number"`
	Country      string  `json:"country"`
	Address      *string `json:"address,omitempty"`
}

// =========== Organizer profile response ==============
type OrganizerProfilesResponse struct {
	ID            string                  `json:"id"`
	UserID        string                  `json:"user_id"`
	Status        OrganizerStatusResp     `json:"status,omitempty"`
	IsVerified    bool                    `json:"is_verified"`
	Email         string                  `json:"email"`
	Name          string                  `json:"name"`
	PhotoProfile  *string                 `json:"photo_profile,omitempty"`
	PhoneNumber   string                  `json:"phone_number"`
	Country       string                  `json:"country"`
	Address       *string                 `json:"address,omitempty"`
	SocialMedia   OrganizerSocialMediaReq `json:"social_media"`
	CompanyDetail OrganizerCompDetailRes  `json:"company_detail"`
	CreatedAt     int64                   `json:"created_at"`
	UpdatedAt     int64                   `json:"updated_at"`
}
