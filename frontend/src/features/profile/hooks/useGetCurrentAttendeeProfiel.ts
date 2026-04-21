import { GetCurrentAttendeeProfileApi } from "../api/profileApi"
import { useQuery } from "@tanstack/react-query";

export const useGetCurrentAttendeeProfile = (enabled: boolean) => {
    return useQuery({
        queryKey: ["attendee-profile"],
        queryFn: async () => {
            const res = await GetCurrentAttendeeProfileApi();
            return res.data;
        },
        staleTime: 1000 * 60 * 5,
        enabled,
    });
}