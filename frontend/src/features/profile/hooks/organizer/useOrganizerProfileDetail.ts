import { useQuery } from "@tanstack/react-query";
import { organizerKeys } from "../../utils/cacheKey";
import { GetOrganizerDetailProfileApi } from "../../api/profileApi";

export const useOrganizerProfileDetail = (id: string) => {
    return useQuery({
        queryKey: organizerKeys.detail(id),
        queryFn: async() => {
            const res = await GetOrganizerDetailProfileApi(id);
            return res.data;
        },
        enabled: !!id,
        staleTime: 1000 * 60 * 5,
    });
};