import { useQuery } from "@tanstack/react-query";
import type { FilterEventsRequest } from "../types/eventRequest";
import { getEventsActiveApi } from "../api/eventsApi";
import { eventsKeys } from "../../../utils/cacheKey";

export const useGetEventsActive = (params: FilterEventsRequest) => {
    return useQuery({
        queryKey: eventsKeys.list(params),
        queryFn: async () => {
            const res = await getEventsActiveApi(params)
            return res.data;
        },
        staleTime: 1000 * 60 * 5,
    })
}