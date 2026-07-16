import { useQuery } from "@tanstack/react-query";
import type { FilterEventsRequest } from "../types/eventRequest";
import { getEventsActiveApi } from "../api/eventsApi";
import { eventsKeys } from "../../../utils/cacheKey";

export const useGetEventsForMap = (
    params: Omit<FilterEventsRequest, "pagination"> & { limit?: number },
    enabled: boolean = true
) => {
    return useQuery({
        queryKey: eventsKeys.list({ ...params, type: "map" }),
        queryFn: async () => {
            const res = await getEventsActiveApi({
                ...params,
                pagination: {
                    page: 1,
                    limit: params.limit || 100,
                },
            });
            return res.data;
        },
        staleTime: 1000 * 60 * 5,
        refetchInterval: 5000,
        enabled,
    });
};