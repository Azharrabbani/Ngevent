import { useMutation, useQueryClient } from "@tanstack/react-query"
import { updateAttendeeProfile } from "../../api/profileApi";
import toast from "react-hot-toast";
import { attendeeKeys } from "../../utils/cacheKey";


export const useUpdateAttendeeProfile = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: updateAttendeeProfile,
        onSuccess: (success) => {
            toast.success(success.data);
            
            queryClient.invalidateQueries({ queryKey: attendeeKeys.lists() });
            queryClient.invalidateQueries({ queryKey: attendeeKeys.details() });
            queryClient.invalidateQueries({ queryKey: attendeeKeys.me() });
        },
        onMutate: async (newData) => {
            await queryClient.cancelQueries({queryKey: attendeeKeys.lists()});

            const previous = queryClient.getQueryData(attendeeKeys.lists());

            queryClient.setQueryData(attendeeKeys.lists(), (old: any) => ({
                ...old,
                ...newData,
            }));

            return { previous }
        },
        onError: (err: any, _, context) => {
            toast.error(err?.response?.data?.message || "Update failed");
            queryClient.setQueryData(attendeeKeys.lists(), context?.previous);
        },
    });
};