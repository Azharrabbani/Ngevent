export interface EventsResponse {
    id: string;
    eo_profile: EoProfile;
    event: EventDetail;
    event_address: EventAddress;
    date: number;
    created_at: number;
    updated_at: number;
    deleted_at?: number;
} 

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