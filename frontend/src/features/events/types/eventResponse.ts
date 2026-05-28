export interface EventsResponse {
    id: string;
    eo_profile: EoProfile;
    event: EventDetail;
    event_address: EventAddress;
    start_time: number;
    end_time: number;
    created_at: number;
    updated_at: number;
    deleted_at?: number;
    update_request_id?: string;
}

export interface UpdateEventResponse {
    id: string;
    event_id: string;
    event_title: string;
    eo_profile: EoProfile;
    updated_details: UpdatedDetails;
    updated_address: EventAddress;
    updated_categories: EventCategories[];
    created_at: number;
    updated_at: number;
    deleted_at?: number;
}
type UpdatedDetails = {
    banner: string;
    status: string;
    description: string;
    start_time: number;
    end_time: number;
    rejected_reason?: string;
    reviewed_by?: Reviewer;
    reviewed_at?: number;
};

type EoProfile = {
    id: string;
    is_verified: boolean;
    email: string;
    name: string;
    photo_profile?: string;
    phone_number: string;
};

type EventDetail = {
    banner?: string;
    name: string;
    categories: EventCategories[];
    tickets: Tickets[];
    slug: string;
    status: string;
    description: string;
    rejected_reason?: string;
    reviewed_by?: Reviewer;
    reviewed_at?: number;
};

type Reviewer = {
    id?: string;
    email?: string;
};

type EventAddress = {
    address: string;
    city: string;
    country: string;
    detail_address: string;
    coordinates: Coordinates;
};

type Coordinates = {
    lat: number;
    lon: number;
};

type EventCategories = {
    id: string;
    name: string;
};

type Tickets = {
    id: string;
    name: string;
    price: string;
    quantity: number;
    ticket_type: string;
};