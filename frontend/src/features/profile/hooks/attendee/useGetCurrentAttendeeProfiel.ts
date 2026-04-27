import { useQuery } from "@tanstack/react-query";
import { GetCurrentAttendeeProfileApi } from "../../api/profileApi";
import { attendeeKeys } from "../../utils/cacheKey";

export const useGetCurrentAttendeeProfile = (enabled: boolean) => {
    return useQuery({
        queryKey: attendeeKeys.me(),
        queryFn: async () => {
            const res = await GetCurrentAttendeeProfileApi();
            return res.data;
        },
        staleTime: 1000 * 60 * 5,
        enabled,
    });
}