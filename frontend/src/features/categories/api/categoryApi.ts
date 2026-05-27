import api from "../../../lib/api"
import type { PaginatedResponse, successResponse } from "../../../types/apiResponse";
import type { FilterCategoryRequest } from "../types/categoryRequest";
import type { categoriesPaginatedResp, categoriesResp } from "../types/categoryResponse"

export const listCategories = async () => {
    const res = await api.get<successResponse<categoriesResp[]>>("/category");
    return res.data
}


export const listCategoriesPaginatedApi = async (request: FilterCategoryRequest) => {
    const res = await api.get<PaginatedResponse<categoriesPaginatedResp>>(
        "/category/list",
        {
            params: {
                name: request.name,
                sort: request.pagination?.sort,
                page: request.pagination?.page,
                limit: request.pagination?.limit
            }
        }
    );
    return res.data
}
