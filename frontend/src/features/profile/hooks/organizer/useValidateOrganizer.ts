import { useMutation, useQueryClient } from "@tanstack/react-query";
import { ValidateOrganizerUpdateApi } from "../../api/profileApi";
import toast from "react-hot-toast";
import { organizerKeys, organizerUpdateKeys } from "../../../../utils/cacheKey";

export const useValidateOrganizerUpdate = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: ({ id, payload }: { id: string; payload: { status: string; reason: string } }) =>
            ValidateOrganizerUpdateApi(id, payload),

        onSuccess: () => {
            toast.success("Review submitted successfully");

            queryClient.invalidateQueries({
                queryKey: organizerKeys.lists(),
                exact: false,
            });

            queryClient.invalidateQueries({
                queryKey: organizerKeys.details(),
                exact: false,
            });

            queryClient.invalidateQueries({
                queryKey: organizerUpdateKeys.all,
                exact: false,
            });
        },

        onError: (err: any) => {
            toast.error(err?.response?.data?.error || "Review failed");
        },
    });
};