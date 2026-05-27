import { useMutation, useQueryClient } from "@tanstack/react-query";
import { ReviewUpdatedEventApi } from "../api/eventsApi";
import { eventsKeys, updateEventKeys } from "../../../utils/cacheKey";
import toast from "react-hot-toast";

export function useReviewUpdatedEvent() {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: ({ id, ...payload }: { id: string; status: "approved" | "rejected"; reason?: string }) =>
            ReviewUpdatedEventApi(id, payload),
        onSuccess: () => {
            toast.success("Event update successfully reviewed")
            queryClient.invalidateQueries({ queryKey: eventsKeys.all });
            queryClient.invalidateQueries({ queryKey: updateEventKeys.all });
        },
        onError: (error: any) => {
            toast.error(error?.message || "Failed to review event update")
        }
    });
}