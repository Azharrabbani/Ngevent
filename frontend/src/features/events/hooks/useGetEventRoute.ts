import { useQuery } from "@tanstack/react-query";
import type { UserLatLonRequest } from "../types/eventRequest";
import { eventsKeys } from "../../../utils/cacheKey";
import { GetEventRouteApi } from "../api/eventsApi";

export const useGetEventRoute = (id: string, params: UserLatLonRequest, enabled: boolean = false) => {
    return useQuery({
        queryKey: [...eventsKeys.route(id), params],
        queryFn: async () => {
            const res = await GetEventRouteApi(id, params)
            return res.data;
        },
        staleTime: 1000 * 60 * 5,
        enabled
    })
}