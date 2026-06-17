import type { PaginationParams } from "../../../types/apiRequest";

export interface FilterEventsRequest {
    search?: string;
    sort?: string;
    date?: string;
    title?: string;
    category?: number[];
    status?: string;
    start_time?: number;
    month?: number;
    year?: number;
    lat?: number;
    lon?: number;
    location?: string;
    with_deleted?: boolean;
    pagination?: PaginationParams;
};

export interface FilterOrganizerEventsRequest {
    title?: string;
    status: string;
    pagination?: PaginationParams;
};

export interface UserLatLonRequest {
    lat?: number;
    lon?: number;
}

export interface FilterUpdatedEventsRequest {
    title?: string;
    search?: string;
    sort?: string;
    date?: string;
    month?: number;
    year?: number;
    status?: string;
    pagination?: PaginationParams;
}
