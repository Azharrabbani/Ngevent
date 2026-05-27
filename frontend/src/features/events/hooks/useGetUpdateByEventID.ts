import { useQuery } from "@tanstack/react-query";
import { GetUpdateByEventIDApi } from "../api/eventsApi";
import { updateEventKeys } from "../../../utils/cacheKey";

export const useGetUpdateByEventID = (eventID: string, status: string) => {
    return useQuery({
        queryKey: updateEventKeys.detail(eventID, status),
        queryFn: async () => {
            const res = await GetUpdateByEventIDApi(eventID, status);
            return res.data;
        },
        staleTime: 1000 * 60 * 5,
        retry: false,
        throwOnError: false,
    });
};