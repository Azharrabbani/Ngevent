export interface RegisterRequest{
    email: string
    password: string
    confirm_password: string
};

export interface LoginRequest {
    email: string
    password: string
};

export interface ForgetPasswordRequest {
    email: string
};

export interface ResetPasswordRequest {
    new_password: string
    confirm_password: string
};

export interface VerifyEmailRequest {
    email: string
    otp: string
};

export interface ResendOtpRequest {
    email: string
};
