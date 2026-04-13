import type { ForgetPasswordRequest, LoginRequest, RegisterRequest, ResendOtpRequest, ResetPasswordRequest, SelectRoleRequest, VerifyEmailRequest } from "../types/authRequest";
import type { LoginResponse, RegisterResponse } from "../types/authResponse";
import type { successResponse } from "../../../types/apiResponse";
import { api } from "../../../lib/api";

export const registerApi = async(payload: RegisterRequest) => {
    const res = await api.post<successResponse<RegisterResponse>>("/user/register", payload)
    return res.data
}

export const loginApi = async(payload: LoginRequest) => {
    const res = await api.post<successResponse<LoginResponse>>("/login", payload)
    return res.data
}

export const forgetPasswordApi = async(payload: ForgetPasswordRequest) => {
    const res = await api.post<successResponse<string>>("/forgot-password", payload)
    return res.data
}

export const resetPasswordApi = async(payload: ResetPasswordRequest, token: string) => {
    const res = await api.put<successResponse<string>>(`/reset-password/${token}`, payload)
    return res.data
}

export const verifyEmailApi = async(payload: VerifyEmailRequest) => {
    const res = await api.put<successResponse<string>>("/verify-email", payload)
    return res.data
}

export const resendOtpApi = async(payload: ResendOtpRequest) => {
    const res = await api.post<successResponse<string>>("/resend-otp", payload);
    return res.data;
};

export const selectRoleApi = async(payload: SelectRoleRequest) => {
    const res = await api.put<successResponse<string>>("/user/role", payload);
    return res.data;
};