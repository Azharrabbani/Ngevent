import type { PaginationParams } from "../../../types/apiRequest";

export interface FilterEventsRequest {
    search?: string;
    sort?: string;
    date?: string;
    get_update?: boolean;
    title?: string;
    category?: number[];
    status?: string;
    start_time?: number;
    location?: string;
    with_deleted?: boolean;
    pagination?: PaginationParams;
};