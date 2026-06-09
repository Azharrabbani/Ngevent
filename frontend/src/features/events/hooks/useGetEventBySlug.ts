import { useQuery } from "@tanstack/react-query";
import type { UserLatLonRequest } from "../types/eventRequest";
import { GetEventBySlug } from "../api/eventsApi";
import { eventsKeys } from "../../../utils/cacheKey";

export const useGetEventBySlug = (
    slug: string,
    params: UserLatLonRequest
) => {
    return useQuery({
        queryKey: [
            ...eventsKeys.detail(slug),
            params.lat,
            params.lon,
        ],

        queryFn: async () => {
            const res = await GetEventBySlug(
                slug,
                params
            );

            return res.data;
        },

        staleTime: 1000 * 60 * 5,
    });
};