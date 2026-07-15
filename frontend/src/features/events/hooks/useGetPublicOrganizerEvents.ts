import { useQuery } from "@tanstack/react-query";
import type { FilterOrganizerEventsRequest } from "../types/eventRequest";
import { eventsKeys } from "../../../utils/cacheKey";
import { GetPublicOrganizerEventsApi } from "../api/eventsApi";

export const useGetPublicOrganizerEvents = (
    id: string,
    params: FilterOrganizerEventsRequest
) => {
    return useQuery({
        queryKey: eventsKeys.list({ ...params, id }),
        queryFn: async () => {
            const res = await GetPublicOrganizerEventsApi(id, params);
            return res.data;
        },
        enabled: !!id,
        staleTime: 1000 * 60 * 5,
        refetchInterval: 5000,
    });
};