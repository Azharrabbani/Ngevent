import { useQuery } from "@tanstack/react-query"
import { eventsKeys } from "../../../../utils/cacheKey"
import type { FilterEventsRequest } from "../../types/organizerRequest"
import { GetEvents } from "../api/eventOrganizerApi"

export const useGetEvents = (params: FilterEventsRequest) => {
    return useQuery({
        queryKey: eventsKeys.list(params),
        queryFn: async () => {
            const res = await GetEvents(params);
            return res.data;
        },
        staleTime: 1000 * 60 * 5,
    });
};