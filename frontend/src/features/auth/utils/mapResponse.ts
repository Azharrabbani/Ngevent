import type { AuthModel } from "../types/authModel";
import type { LoginResponse } from "../types/loginResponse";

export const mapLoginResponse = (res: LoginResponse): AuthModel => {
    return {
        id: res.data.id,
        email: res.data.email,
        role: res.data.role,
        token: res.data["ngevent-token"],
        refreshToken: res.data["ngevent-ref-token"]
    }
}