import { useInfiniteQuery } from "@tanstack/react-query";
import type { FilterOrganizerReq } from "../../types/profileRequest";
import { eventsPublicKeys } from "../../../../utils/cacheKey";
import { GetPublicOrganizersApi } from "../../api/profileApi";

export const useGetPublicOrganizers = (
    params: Omit<FilterOrganizerReq, "pagination"> & { limit?: number }
) => {
    return useInfiniteQuery({
        queryKey: eventsPublicKeys.list({ ...params, type: "infinite" }),
        queryFn: async ({ pageParam = 1 }: { pageParam: number }) => {
            const res = await GetPublicOrganizersApi({
                ...params,
                pagination: {
                    page: pageParam,
                    limit: params.limit || 12,
                },
            });
            return res.data;
        },
        initialPageParam: 1,
        getNextPageParam: (lastPage) => {
            const nextPage = lastPage.page + 1;
            return nextPage <= lastPage.total_pages ? nextPage : undefined;
        },
        staleTime: 1000 * 60 * 10,
    });
};