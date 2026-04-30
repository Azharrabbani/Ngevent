import api from "../../../lib/api"
import type { successResponse } from "../../../types/apiResponse"
import type { categoriesResp } from "../types/categoryResponse"

export const listCategories = async () => {
    const res = await api.get<successResponse<categoriesResp[]>>("/category");
    return res.data
}