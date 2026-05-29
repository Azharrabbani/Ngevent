import type { PaginationParams } from "../../../types/apiRequest";

export interface UserFilterRequest {
    role?: string;
    isVerified?: boolean;
    email?: string;
    pagination?: PaginationParams;
};