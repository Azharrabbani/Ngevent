import axios from "axios";
import type { LoginRequest } from "../types/loginRequest";
import type { LoginResponse } from "../types/loginResponse";

const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL,
    withCredentials: true,
})

export const loginApi = async(payload: LoginRequest) => {
    const res = await api.post<LoginResponse>("/login", payload)
    return res.data
}