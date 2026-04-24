import { useQuery } from "@tanstack/react-query"
import { ListUsersApi } from "../api/profileApi"
import type { UserQueryParams } from "../types/profileRequest";

export const useListUsers = (params: UserQueryParams) => {
    return useQuery ({
        queryKey: ["all-users", params],
        queryFn: async () => {
            const res = await ListUsersApi(params);
            return res.data;
        },
        staleTime: 1000 * 60 * 5,
    })
}