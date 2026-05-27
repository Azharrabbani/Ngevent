import { useQuery } from "@tanstack/react-query"
import { categoriesKeys } from "../../../utils/cacheKey"
import type { FilterCategoryRequest } from "../types/categoryRequest"
import { listCategoriesPaginatedApi } from "../api/categoryApi"

export const useListCategoriesPaginated = (params: FilterCategoryRequest) => {
    return useQuery({
        queryKey: categoriesKeys.list(params),
        queryFn: async () => {
            const res = await listCategoriesPaginatedApi(params);
            return res.data;
        },
        staleTime: 1000 * 60 * 5,
    })
}