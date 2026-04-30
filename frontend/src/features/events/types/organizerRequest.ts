import type { PaginationParams } from "../../../types/apiRequest";

export interface FilterEventsRequest {
    title?: string;
    category?: number[];
    status?: string;
    date?: number;
    location?: string;
    pagination?: PaginationParams;
};