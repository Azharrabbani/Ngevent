export interface UserResponse {
    id: string;
    email: string
    role: string | null
    has_profile: boolean | null
    is_verified: boolean
    created_at: number
    updated_at: number
    deleted_at: number | null
};