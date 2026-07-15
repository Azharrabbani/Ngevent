// useCreateOrganizerProfile.ts
import { useMutation, useQueryClient } from "@tanstack/react-query"
import type { CreateOrganizerProfileReq } from "../../types/profileRequest"
import { CreateOrganizerProfileApi } from "../../api/profileApi"
import toast from "react-hot-toast"
import { organizerKeys } from "../../../../utils/cacheKey"

export const useCreateOrganizerProfile = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (payload: CreateOrganizerProfileReq) => CreateOrganizerProfileApi(payload),
        onSuccess: (res) => {
            toast.success(res.data)

            queryClient.invalidateQueries({ queryKey: organizerKeys.lists() });
            queryClient.invalidateQueries({ queryKey: organizerKeys.details() });
            queryClient.invalidateQueries({ queryKey: organizerKeys.me() });
        },
        onError: (err: any) => {
            const msg = err?.response?.data?.error
            console.log("msg: ", msg)
            if (!Array.isArray(msg)) {
                toast.error(msg || "Failed to create profile");
            }
        }
    })
}