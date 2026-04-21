import { useMutation, useQueryClient } from "@tanstack/react-query"
import { updateAttendeeProfile } from "../api/profileApi";
import toast from "react-hot-toast";


export const useUpdateAttendeeProfile = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: updateAttendeeProfile,
        onSuccess: (success) => {
            toast.success(success.data);
            queryClient.invalidateQueries({
                queryKey: ["attendee-profile"],
            });
        },
        onMutate: async (newData) => {
            await queryClient.cancelQueries({queryKey: ["attendee-profile"]});

            const previous = queryClient.getQueryData(["attendee-profile"]);

            queryClient.setQueryData(["attendee-profile"], (old: any) => ({
                ...old,
                ...newData,
            }));

            return { previous }
        },
        onError: (err: any, _, context) => {
            toast.error(err?.response?.data?.message || "Update failed");
            queryClient.setQueryData(["attendee-profile"], context?.previous);
        },
    });
};