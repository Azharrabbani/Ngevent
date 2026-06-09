import { useMutation, useQueryClient } from "@tanstack/react-query"
import toast from "react-hot-toast"
import { CancelEventApi } from "../api/eventsApi"
import { eventsKeys, eventsPublicKeys } from "../../../utils/cacheKey"

export const useCancelEvent = () => {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: (id: string) => CancelEventApi(id),
        onSuccess: (success) => {
            toast.success(success.data)
            // Invalidate list queries
            queryClient.invalidateQueries({
                queryKey: eventsKeys.all,
                refetchType: "active"
            });

            queryClient.invalidateQueries({
                queryKey: eventsPublicKeys.all,
                refetchType: "active"
            });
        },
        onError: (err: any) => {
            toast.error(err.response?.data?.error || "Failed to cancel event")
        }
    })
}