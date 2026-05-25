import { useQuery } from "@tanstack/react-query";
import type { FilterUpdatedEventsRequest } from "../types/eventRequest";
import { GetAllUpdatedEventsApi } from "../api/eventsApi";
import { updateEventKeys } from "../../../utils/cacheKey";

export function useGetAllUpdatedEvents(
    params: FilterUpdatedEventsRequest,
    enabled: boolean = true
) {
    return useQuery({
        queryKey: updateEventKeys.list(params),
        queryFn: async () => {
            const res = await GetAllUpdatedEventsApi(params);
            return res.data;
        },
        retry: false,
        staleTime: 1000 * 60 * 5,
        enabled,
    });
}