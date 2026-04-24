export interface AttendeeResponse {
    id: string;
    user_id: string;
    email: string;
    name: string;
    username: string;
    photo_profile: string;
    phone_number: string;
    country: string;
    address: string;
};

export interface OrganizerResponse {
    id: string;
    user_id: string;
    status: OrganizerStatus;
    is_verified: boolean;
    email: string;
    name : string;
    photo_profile: string;
    phone_number: string;
    country: string;
    address: string;
    social_media: OrganizerSocialMedia;
    company_detail: OrganizerCompanyDetail;
};

type OrganizerStatus = {
    status: string;
    rejected_reason: string;
    reviewed_by: string;
    reviewed_at: string;
};

type OrganizerSocialMedia = {
    email: string;
    instagram: string;
};

type OrganizerCompanyDetail = {
    description: string;
    npwp: string;
    npwp_file: string;
    nib: string;  
    nib_file: string;
};

export interface User {
    id: string;
    email: string;
    role: string;
    is_verified: boolean;
    created_at: number;
    updated_at: number;
    deleted_at: number;
};
