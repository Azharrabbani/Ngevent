import type { PaginationParams } from "../../../types/apiRequest";

// Attendee data
export interface CreateAttendeeProfileReq {
    photo: File;
    name: string;
    username: string;
    phonenumber: string;
    iso: string;
    address: string;
};

export interface UpdateAttendeeProfileReq {
    name: string;
    username: string;
    phone_number: string;
    iso: string;
    address: string;
};

export interface FilterAttendeeReq {
    filter?: string;
    pagination?: PaginationParams;
};


// Organizer data
export interface CreateOrganizerProfileReq {
    photo: File;
    name: string;
    phonenumber: string;
    iso: string;
    address: string;
    description: string;
    npwp: string;
    npwpFile: File;
    nib: string;
    nibFile: File;
}

export interface UpdateOrganizerProfileReq {
    name: string;
    phonenumber: string;
    iso: string;
    address: string;
    description: string;
    email: string;
    instagram: string;
    npwp: string;
    npwpFile: File;
    nib: string;
    nibFile: File;
}


export interface UpdatePhotoReq {
    photo: File;
};

export interface UserQueryParams {
    role?: string;
    is_verified?: boolean;
    email?: string;
    pagination?: PaginationParams
};
