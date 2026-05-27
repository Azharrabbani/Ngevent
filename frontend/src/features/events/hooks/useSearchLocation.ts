import { useQuery } from "@tanstack/react-query"
import { locationKeys } from "../../../utils/cacheKey"
import type { locationReq } from "../types/locationRequest"
import { SearchLocationApi } from "../api/eventsApi"

export const useSearchLocation = (
    payload: locationReq,
    options?: { enabled?: boolean }
) => {
    return useQuery({
        queryKey: locationKeys.list(payload),
        queryFn: async () => {
            const res = await SearchLocationApi(payload);
            return res.data;
        },
        enabled: options?.enabled,
        staleTime: 1000 * 60 * 5,
    });
};