import { GetCurrentOrganizerProfileApi } from "../../api/profileApi";
import { useQuery } from "@tanstack/react-query";
import { organizerKeys } from "../../utils/cacheKey";

export const useGetCurrentOrganizerProfile = (enabled: boolean) => {
    return useQuery({
        queryKey: organizerKeys.me(),
        queryFn: async () => {
            const res = await GetCurrentOrganizerProfileApi();
            return res.data;
        },
        enabled,
    });
}