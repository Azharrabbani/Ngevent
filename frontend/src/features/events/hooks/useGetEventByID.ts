import { useQuery } from "@tanstack/react-query"
import { eventsKeys } from "../../../utils/cacheKey"
import { GetEventByID } from "../api/eventsApi";

export const useGetEventByID = (id: string) => {
    return useQuery({
        queryKey: eventsKeys.detail(id),
        queryFn: async () => {
            const res = await GetEventByID(id);
            return res.data;
        },
        enabled: !!id,
        staleTime: 1000 * 60 * 5,
    });
};