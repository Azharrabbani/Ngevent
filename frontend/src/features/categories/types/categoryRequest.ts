import type { PaginationParams } from "../../../types/apiRequest";

export interface FilterCategoryRequest {
    name?: string;
    pagination?: PaginationParams;
};

export interface CreateCategoryRequest {
    name: string
}