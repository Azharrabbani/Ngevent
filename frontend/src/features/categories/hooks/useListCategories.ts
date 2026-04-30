import { useQuery } from "@tanstack/react-query"
import { categoriesKeys } from "../../../utils/cacheKey"
import { listCategories } from "../api/categoryApi"

export const useListCategories = () => {
    return useQuery({
        queryKey: categoriesKeys.lists(),
        queryFn: async () => {
            const res = await listCategories();
            return res.data;
        },
        staleTime: 1000 * 60 * 5,
    });
};