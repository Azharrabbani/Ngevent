export interface LoginResponse {  
    id: string
    email: string
    role: string
    "ngevent-token": string
    "ngevent-ref-token": string
}

export interface RegisterResponse {
    id: string
    email: string
    role: string
    is_verified: string
    created_at: number
    updated_at: number
    deleted_at: number
}

