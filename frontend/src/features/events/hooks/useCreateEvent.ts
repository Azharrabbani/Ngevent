import { useMutation, useQueryClient } from "@tanstack/react-query"
import toast from "react-hot-toast";
import { CreateEventApi } from "../api/eventsApi";
import { eventsKeys } from "../../../utils/cacheKey";

type ValidationError = { field: string; message: string };

export const useCreateEvent = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: ({ payload, banner }: any) =>
            CreateEventApi(payload, banner),

        onError: (err: any) => {
            const validationErrors: ValidationError[] | undefined =
                err?.response?.data?.error;

            if (Array.isArray(validationErrors)) {
                validationErrors.forEach(({ message }) => toast.error(message));
                return;
            }

            toast.error(err?.response?.data?.error || "Failed to create event");
        },

        onSuccess: (success) => {
            toast.success(success.data);

            queryClient.invalidateQueries({
                queryKey: eventsKeys.lists(),
                exact: false,
                refetchType: "active",
            });
        },
    });
};