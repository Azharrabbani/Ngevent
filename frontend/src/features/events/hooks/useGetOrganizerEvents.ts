import { useQuery } from "@tanstack/react-query"
import { eventsKeys } from "../../../utils/cacheKey"
import type { FilterEventsRequest } from "../types/eventRequest"
import { GetOrganizerEventsApi } from "../api/eventsApi";

export const useGetOrganizerEvents = (params: FilterEventsRequest) => {
    return useQuery({
        queryKey: eventsKeys.list(params),
        queryFn: async () => {
            const res = await GetOrganizerEventsApi(params);
            return res.data;
        },
        staleTime: 1000 * 60 * 5,
        refetchInterval: 5000,
    });
};