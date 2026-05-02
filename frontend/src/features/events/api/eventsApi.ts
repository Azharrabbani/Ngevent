import api from "../../../lib/api"
import type { successResponse } from "../../../types/apiResponse"
import type { locationReq } from "../types/locationRequest"
import type { locationResp } from "../types/locationResponse"

export const SearchLocationApi = async (payload: locationReq) => {
    const res = await api.get<successResponse<locationResp[]>>("/location/", {
        params: {
            query: payload.query
        },
    });
    
    return res.data;
}