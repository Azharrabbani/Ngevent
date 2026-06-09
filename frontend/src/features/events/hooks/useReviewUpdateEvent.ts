import { useMutation, useQueryClient } from "@tanstack/react-query";
import { ReviewUpdatedEventApi } from "../api/eventsApi";
import { eventsKeys, eventsPublicKeys, updateEventKeys } from "../../../utils/cacheKey";
import toast from "react-hot-toast";

export function useReviewUpdatedEvent() {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: ({ id, ...payload }: { id: string; status: "approved" | "rejected"; reason?: string }) =>
            ReviewUpdatedEventApi(id, payload),
        onSuccess: () => {
            toast.success("Event update successfully reviewed")
            queryClient.invalidateQueries({
                queryKey: eventsKeys.all,
                refetchType: "active",
            });
            queryClient.invalidateQueries({
                queryKey: updateEventKeys.all,
                refetchType: "active",
            });
            queryClient.invalidateQueries({
                queryKey: eventsPublicKeys.all,
                refetchType: "active",
            });
        },
        onError: (error: any) => {
            toast.error(error?.message || "Failed to review event update")
        }
    });
}