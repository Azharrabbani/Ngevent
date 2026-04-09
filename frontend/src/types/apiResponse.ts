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