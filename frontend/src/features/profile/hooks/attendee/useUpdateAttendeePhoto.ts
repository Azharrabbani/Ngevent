import toast from "react-hot-toast";
import { updateAttendeePhotoApi } from "../../api/profileApi"
import { useMutation, useQueryClient } from "@tanstack/react-query"
import { attendeeKeys } from "../../utils/cacheKey";

export const useUpdateAttendeePhoto = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: updateAttendeePhotoApi,
        onMutate: async (newData) => {
            await queryClient.cancelQueries({queryKey: attendeeKeys.lists()});

            const previous = queryClient.getQueryData(attendeeKeys.lists());

            const previewUrl = URL.createObjectURL(newData.photo);

            queryClient.setQueryData(attendeeKeys.lists(), (old: any) => ({
                ...old,
                photo_profile: previewUrl,
            }));

            return { previous, previewUrl };
        },
        onError: (err: any, _, context) => {
            toast.error(err?.response?.data?.message || "Update failed")
            queryClient.setQueryData(attendeeKeys.lists(), context?.previous);

            if (context?.previewUrl) {
                URL.revokeObjectURL(context.previewUrl);
            }
        },
        onSuccess: (_, __, context) => {
            if (context?.previewUrl) {
                URL.revokeObjectURL(context.previewUrl);
            }

            queryClient.invalidateQueries({ queryKey: attendeeKeys.lists() });
            queryClient.invalidateQueries({ queryKey: attendeeKeys.details() });
            queryClient.invalidateQueries({ queryKey: attendeeKeys.me() });
        },
    });
};