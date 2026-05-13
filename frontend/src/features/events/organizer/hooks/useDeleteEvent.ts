import { useMutation, useQueryClient } from "@tanstack/react-query"
import { DeleteEventApi } from "../api/eventOrganizerApi"
import toast from "react-hot-toast"
import { eventsKeys } from "../../../../utils/cacheKey"

export const useDeleteEvent = () => {
    const queryClient = useQueryClient()

    return useMutation({
        mutationFn: (id: string) => DeleteEventApi(id),
        onSuccess: (success, variable) => {
            toast.success(success.data)
            // Invalidate list queries
            queryClient.invalidateQueries({
                queryKey: eventsKeys.lists(),
                exact: false,
                refetchType: "active"
            });

            // Invalidate the specific event detail
            queryClient.invalidateQueries({
                queryKey: eventsKeys.detail(variable),
                refetchType: "active"
            });
        },
        onError: (err: any) => {
            toast.error(err.response?.data?.error || "Failed to delete event")
        }
    })
}