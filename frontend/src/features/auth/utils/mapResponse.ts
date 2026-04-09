import type { successResponse } from "../../../types/apiResponse";
import type { AuthModel } from "../types/authModel";

export const mapLoginResponse = (res: successResponse<any>): AuthModel => {
    return {
        id: res.data.id,
        email: res.data.email,
        role: res.data.role,
        token: res.data["ngevent-token"],
        refreshToken: res.data["ngevent-ref-token"]
    }
}