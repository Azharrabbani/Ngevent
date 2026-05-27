import type { PaginationParams } from "../../../types/apiRequest";

export interface FilterEventsRequest {
    search?: string;
    sort?: string;
    date?: string;
    title?: string;
    category?: number[];
    status?: string;
    start_time?: number;
    location?: string;
    with_deleted?: boolean;
    pagination?: PaginationParams;
};

export interface FilterUpdatedEventsRequest {
    title?: string;
    search?: string;   
    sort?: string;     
    date?: string;     
    status?: string;
    pagination?: PaginationParams;
}
