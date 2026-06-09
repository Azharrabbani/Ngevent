export interface EventCategory {
    id: number;
    name: string;
}

export interface EoProfile {
    id: string;
    is_verified: boolean;
    status: string;
    email: string;
    name: string;
    photo_profile: string;
    phone_number: string;
}

export interface EventInfo {
    banner: string;
    name: string;
    categories: EventCategory[];
    slug: string;
    status: string;
    request_updates: boolean;
    description: string;
}

export interface Coordinates {
    lat: number;
    lon: number;
}

export interface EventAddress {
    address: string;
    city: string;
    country: string;
    detail_address: string;
    coordinates: Coordinates;
}

export interface PathPoint {
    name: string;
    lat: number;
    lon: number;
}

export interface EventDetailResponse {
    id: string;
    eo_profile: EoProfile;
    event: EventInfo;
    event_address: EventAddress;
    start_time: number;
    end_time: number;
    distance?: string;
    path?: PathPoint[];
    created_at: number;
    updated_at: number;
}