import { GetCurrentOrganizerProfileApi } from "../api/profileApi";
import { useQuery } from "@tanstack/react-query";

export const useGetCurrentOrganizerProfile = (enabled: boolean) => {
    return useQuery({
        queryKey: ["organizer-profile"],
        queryFn: async () => {
            const res = await GetCurrentOrganizerProfileApi();
            return res.data;
        },
        enabled,
    });
}