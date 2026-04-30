import { useQuery } from "@tanstack/react-query"
import { GetAttendeeDetailProfileApi } from "../../api/profileApi"
import { attendeeKeys } from "../../../../utils/cacheKey";

export const useAttendeeProfileDetail = (id: string) => {
    return useQuery({
        queryKey: attendeeKeys.detail(id),
        queryFn: async() => {
            const res = await GetAttendeeDetailProfileApi(id);
            return res.data;
        },
        staleTime: 1000 * 60 * 5,
    });
};