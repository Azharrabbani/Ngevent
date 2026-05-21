import { useMutation, useQueryClient } from "@tanstack/react-query"
import toast from "react-hot-toast";
import { CreateEventApi } from "../api/eventsApi";
import { eventsKeys } from "../../../utils/cacheKey";

export const useCreateEvent = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: ({ payload, banner }: any) => 
            CreateEventApi(payload, banner),
        onError: (err: any) => {
            const validationError = err?.response?.data?.error;

            if (!Array.isArray(validationError)) {
                toast.error(err?.response?.data?.error || "Failed create event");
            }
        },
        onSuccess: (success) => {
            toast.success(success.data);

            queryClient.invalidateQueries({ 
                queryKey: eventsKeys.lists(),
                exact: false,
                refetchType: "active"
            });
        }

    })
}