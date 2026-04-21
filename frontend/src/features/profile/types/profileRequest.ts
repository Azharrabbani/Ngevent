// Attendee data
export interface CreateAttendeeProfileReq {
    photo: File
    name: string
    username: string
    phonenumber: string
    iso: string
    address: string
};

export interface UpdateAttendeeProfileReq {
    name: string
    username: string
    phone_number: string
    iso: string
    address: string
};


// Organizer data
export interface CreateOrganizerProfileReq {
    photo: File
    name: string
    phonenumber: string
    iso: string
    address: string
    description: string
    npwp: string
    npwpFile: File
    nib: string
    nibFile: File
}

export interface UpdateOrganizerProfileReq {
    name: string
    phonenumber: string
    iso: string
    address: string
    description: string
    email: string
    instagram: string
    npwp: string
    npwpFile: File
    nib: string
    nibFile: File
}


export interface UpdatePhotoReq {
    photo: File
};
