import { useMutation, useQueryClient } from "@tanstack/react-query"
import { UpdateOrganizerProfileApi } from "../api/profileApi";
import toast from "react-hot-toast";

export const useUpdateOrganizerProfile = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: UpdateOrganizerProfileApi,
        onSuccess: (success) => {
            toast.success(success.data);
            queryClient.invalidateQueries({
                queryKey: ["organizer-profile"],
            });
        },
        onMutate: async (newData) => {
            await queryClient.cancelQueries({queryKey: ["organizer-profile"]});

            const previous = queryClient.getQueryData(["organizer-profile"]);

            queryClient.setQueryData(["organizer-profile"], (old: any) => ({
                ...old,
                ...newData,
            }));

            return { previous };
        },
        onError: (err: any, _, context) => {
            toast.error(err?.response?.data?.error || "Update failed");
            queryClient.setQueryData(["organizer-profile"], context?.previous);
        },
    });
};