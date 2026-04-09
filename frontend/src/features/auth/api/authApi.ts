import type { ForgetPasswordRequest, LoginRequest, ResetPasswordRequest } from "../types/authRequest";
import type { LoginResponse } from "../types/authResponse";
import type { successResponse } from "../../../types/apiResponse";
import { api } from "../../../lib/api";


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
