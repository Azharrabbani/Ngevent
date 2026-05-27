import { useMutation, useQueryClient } from "@tanstack/react-query"
import toast from "react-hot-toast";
import { UpdateEventApi } from "../api/eventsApi";
import { eventsKeys } from "../../../utils/cacheKey";

export const useUpdateEvent = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: ({ id, payload, banner }: { id: string; payload: any; banner: File | null }) =>
            UpdateEventApi(id, payload, banner),
        onError: (err: any) => {
            const validationError = err?.response?.data?.error;

            if (!Array.isArray(validationError)) {
                toast.error(err?.response?.data?.error || "Failed to update event");
            }
        },
        onSuccess: (success, variables) => {
            toast.success(success.data);

            // Invalidate list queries
            queryClient.invalidateQueries({
                queryKey: eventsKeys.lists(),
                exact: false,
                refetchType: "active"
            });

            // Invalidate the specific event detail
            queryClient.invalidateQueries({
                queryKey: eventsKeys.detail(variables.id),
                refetchType: "active"
            });
        }
    });
};
