export interface EventsResponse {
    id: string;
    eo_profile: EoProfile;
    event: EventDetail;
    event_address: EventAddress;
    start_date: number;
    end_date: number;
    start_time: number;
    end_time: number;
    created_at: number;
    updated_at: number;
    submitted_at?: number;
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
    submitted_at: number;
    deleted_at?: number;
}
type UpdatedDetails = {
    banner: string;
    status: string;
    description: string;
    start_date: number;
    end_date: number;
    start_time: number;
    end_time: number;
    rejected_reason?: string;
    reviewed_by?: Reviewer;
    reviewed_at?: number;
};

type EoProfile = {
    id: string;
    is_verified: boolean;
    status: string;
    email: string;
    name: string;
    slug: string;
    photo_profile?: string;
    phone_number: string;
};

type EventDetail = {
    banner?: string;
    name: string;
    categories: EventCategories[];
    slug: string;
    status: string;
    request_updates?: boolean;
    description: string;
    rejected_reason?: string;
    reviewed_by?: Reviewer;
    reviewed_at?: number;
};

export interface RouteRespone {
    event: string;
    distance: string;
    path: PathResponse[];
}

type PathResponse = {
    name: string;
    lat: number;
    lon: number;
}

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