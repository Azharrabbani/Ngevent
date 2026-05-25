import { useMutation, useQueryClient } from "@tanstack/react-query";
import { ReviewEventApi } from "../api/eventsApi";
import { eventsKeys } from "../../../utils/cacheKey";
import toast from "react-hot-toast";

export function useReviewEvent() {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: ({ id, ...payload }: { id: string; status: "active" | "rejected"; reason?: string }) =>
            ReviewEventApi(id, payload),
        onSuccess: () => {
            toast.success("Event successfully reviewed")
            queryClient.invalidateQueries({ queryKey: eventsKeys.all });
        },
        onError: (error: any) => {
            toast.error(error?.message || "Failed to review event")
        }
    });
}