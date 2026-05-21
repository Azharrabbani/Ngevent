import { useQuery } from "@tanstack/react-query";
import { eventsKeys } from "../../../utils/cacheKey";
import type { FilterEventsRequest } from "../types/organizerRequest";
import { GetEventsApi } from "../api/eventsApi";

export const useGetEvents = (params: FilterEventsRequest) => {
    return useQuery({
        queryKey: eventsKeys.list(params),
        queryFn: async () => {
            const res = await GetEventsApi(params);
            return res.data;
        },
        staleTime: 1000 * 60 * 5,
    });
};