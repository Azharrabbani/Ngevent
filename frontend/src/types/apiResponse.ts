export interface successResponse<T> {
    code: number
    status: string
    message: string
    data: T
}

export interface errorResponse<T> {
    code: number
    status: string
    message: string
    error: T
}

export interface PaginatedData<T> {
    limit: number;
    page: number;
    sort: string;
    total_rows: number;
    total_pages: number;
    rows: T[];
};

export type PaginatedResponse<T> = successResponse<PaginatedData<T>>;
