import { useQuery } from "@tanstack/react-query"
import { organizerKeys } from "../../utils/cacheKey"
import type { FilterOrganizerReq } from "../../types/profileRequest"
import { GetOrganizersProfileApi } from "../../api/profileApi"

export const useListOrganizer = (params: FilterOrganizerReq) => {
    return useQuery({
        queryKey: organizerKeys.list(params),
        queryFn: async() => {
            const res = await GetOrganizersProfileApi(params);

            return res.data;
        },
        staleTime: 1000 * 60 * 5,
    });
};