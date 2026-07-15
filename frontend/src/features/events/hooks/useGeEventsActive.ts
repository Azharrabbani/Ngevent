import { useInfiniteQuery } from "@tanstack/react-query";
import type { FilterEventsRequest } from "../types/eventRequest";
import { getEventsActiveApi } from "../api/eventsApi";
import { eventsKeys } from "../../../utils/cacheKey";

export const useGetEventsActive = (
    params: Omit<FilterEventsRequest, "pagination"> & { limit?: number },
    enabled: boolean = true
) => {
    return useInfiniteQuery({
        queryKey: eventsKeys.list({ ...params, type: "infinite" }),
        queryFn: async ({ pageParam = 1 }: { pageParam: number }) => {
            const res = await getEventsActiveApi({
                ...params,
                pagination: {
                    page: pageParam,
                    limit: params.limit || 8,
                }
            });
            return res.data;
        },
        initialPageParam: 1,
        getNextPageParam: (lastPage) => {
            const nextPage = lastPage.page + 1;
            return nextPage <= lastPage.total_pages ? nextPage : undefined;
        },
        staleTime: 1000 * 60 * 5,
        refetchInterval: 5000,
        enabled,
    });
};