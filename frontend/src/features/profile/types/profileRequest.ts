export interface CreateAttendeeProfileReq {
    photo: File
    name: string
    username: string
    phonenumber: string
    iso: string
    address: string
};

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