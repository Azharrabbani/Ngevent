import { useQuery } from "@tanstack/react-query";
import type { FilterOrganizerReq } from "../../types/profileRequest";
import { eventsPublicKeys } from "../../../../utils/cacheKey";
import { GetPublicOrganizersApi } from "../../api/profileApi";

export const useGetPublicOrganizers = (params: FilterOrganizerReq) => {
    return useQuery({
        queryKey: eventsPublicKeys.list(params),
        queryFn: async () => {
            const res = await GetPublicOrganizersApi(params);
            return res.data;
        },
        staleTime: 10 * 60 * 1000,
    })
}