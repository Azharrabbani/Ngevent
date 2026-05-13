import type { PaginationParams } from "../../../types/apiRequest";

export interface FilterEventsRequest {
    title?: string;
    category?: number[];
    status?: string;
    start_time?: number;
    location?: string;
    with_deleted?: boolean;
    pagination?: PaginationParams;
};