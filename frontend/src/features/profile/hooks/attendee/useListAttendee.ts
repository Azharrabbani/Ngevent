import { useQuery } from "@tanstack/react-query"
import type { FilterAttendeeReq } from "../../types/profileRequest"
import { GetAttendeesProfileApi } from "../../api/profileApi"
import { attendeeKeys } from "../../utils/cacheKey";

export const useListAttendee = (params: FilterAttendeeReq) => {
    return useQuery({
        queryKey: attendeeKeys.list(params),
        queryFn: async () => {
            const res = await GetAttendeesProfileApi(params);
            return res.data;
        },
        staleTime: 1000 * 60 * 5,
    });
};