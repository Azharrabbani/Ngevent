import toast from "react-hot-toast";
import { updateAttendeePhotoApi } from "../api/profileApi"
import { useMutation, useQueryClient } from "@tanstack/react-query"

export const useUpdateAttendeePhoto = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: updateAttendeePhotoApi,
        onMutate: async (newData) => {
            await queryClient.cancelQueries({queryKey: ["attendee-profile"]});

            const previous = queryClient.getQueryData(["attendee-profile"]);

            const previewUrl = URL.createObjectURL(newData.photo);

            queryClient.setQueryData(["attendee-profile"], (old: any) => ({
                ...old,
                photo_profile: previewUrl,
            }));

            return { previous, previewUrl };
        },
        onError: (err: any, _, context) => {
            toast.error(err?.response?.data?.message || "Update failed")
            queryClient.setQueryData(["attendee-profile"], context?.previous);

            if (context?.previewUrl) {
                URL.revokeObjectURL(context.previewUrl);
            }
        },
        onSuccess: (_, __, context) => {
            if (context?.previewUrl) {
                URL.revokeObjectURL(context.previewUrl);
            }

            queryClient.invalidateQueries({
                queryKey: ["attendee-profile"],
            });
        },
    });
};