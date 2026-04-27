import { useQuery } from "@tanstack/react-query"
import { GetUpdateOrganizerReqApi } from "../../api/profileApi"
import { organizerUpdateKeys } from "../../utils/cacheKey";

export const useGetOrganizerUpdate = (id: string) => {
    return useQuery({
        queryKey: organizerUpdateKeys.detail(id),
        queryFn: async () => {
            const res = await GetUpdateOrganizerReqApi(id);
            return res.data;
        },
        staleTime: 1000 * 60 * 5,
        enabled: !!id,
        retry: false, 
    });
};