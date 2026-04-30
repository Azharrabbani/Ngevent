import { useMutation, useQueryClient } from "@tanstack/react-query"
import { UpdateOrganizerProfileApi } from "../../api/profileApi";
import toast from "react-hot-toast";
import { organizerKeys } from "../../../../utils/cacheKey";

export const useUpdateOrganizerProfile = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: UpdateOrganizerProfileApi,

        onMutate: async (newData) => {
            await queryClient.cancelQueries({
                queryKey: organizerKeys.me(),
            });

            const previous = queryClient.getQueryData(
                organizerKeys.me()
            );

            queryClient.setQueryData(
                organizerKeys.me(),
                (old: any) => ({
                    ...old,
                    ...newData,
                })
            );

            return { previous };
        },

        onError: (err: any, _, context) => {
            if (err?.response?.data?.error === "NPWP and NIB file required") {
                toast.error("NPWP and NIB file required");
            } else {
                toast.error(err?.response?.data?.error || "Update failed");
            }


            if (context?.previous) {
                queryClient.setQueryData(
                    organizerKeys.me(),
                    context.previous
                );
            }
        },

        onSuccess: (success) => {
            toast.success(success.data);

            queryClient.invalidateQueries({
                queryKey: organizerKeys.me(),
            });

            queryClient.invalidateQueries({
                queryKey: organizerKeys.lists(),
                exact: false,
            });
        },
    });
};