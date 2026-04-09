export interface LoginRequest{
    email: string
    password: string
}

export interface ForgetPasswordRequest{
    email: string
}

export interface ResetPasswordRequest{
    new_password: string
    confirm_password: string
}
