import { useQuery } from "@tanstack/react-query"
import { listUsersApi } from "../api/profileApi"
import type { UserQueryParams } from "../types/profileRequest";
import { userKeys } from "../../../utils/cacheKey";

export const useListUsers = (params: UserQueryParams) => {
    return useQuery({
        queryKey: userKeys.list(params),
        queryFn: async () => {
            const res = await listUsersApi(params);
            return res.data;
        },
        staleTime: 1000 * 60 * 5,
        refetchInterval: 30000,
    })
}