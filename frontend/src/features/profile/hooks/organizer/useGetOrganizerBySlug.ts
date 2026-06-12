import { useQuery } from "@tanstack/react-query"
import { organizerKeys } from "../../../../utils/cacheKey"
import { GetOrganizerBySlugApi } from "../../api/profileApi"

export const useGetOrganizerBySlug = (slug: string) => {
    return useQuery({
        queryKey: organizerKeys.detail(slug),
        queryFn: async () => {
            const res = await GetOrganizerBySlugApi(slug);
            return res.data;
        },
        staleTime: 1000 * 60 * 5,
    })
}