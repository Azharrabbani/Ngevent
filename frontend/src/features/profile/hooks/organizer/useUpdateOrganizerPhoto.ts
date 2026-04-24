import { useMutation, useQueryClient } from "@tanstack/react-query"

import toast from "react-hot-toast";
import { UpdateOrganizerPhotoApi } from "../../api/profileApi";

export const useUpdateOrganizerPhoto = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: UpdateOrganizerPhotoApi,
        onMutate: async (newData) => {
            await queryClient.cancelQueries({queryKey: ["organizer-profile"]});

            const previous = queryClient.getQueryData(["organizer-profile"]);

            const previewUrl = URL.createObjectURL(newData.photo);

            queryClient.setQueryData(["organizer-profile"], (old: any) => ({
                ...old,
                photo_profile: previewUrl,
            }));

            return { previous, previewUrl };
        },
        onError: (err: any, _, context) => {
            toast.error(err?.response?.data?.message || "Update failed")
            queryClient.setQueryData(["organizer-profile"], context?.previous);

            if (context?.previewUrl) {
                URL.revokeObjectURL(context.previewUrl);
            }
        },
        onSuccess: (_, __, context) => {
            if (context?.previewUrl) {
                URL.revokeObjectURL(context.previewUrl);
            }

            queryClient.invalidateQueries({
                queryKey: ["organizer-profile"],
            });
        },
    });
}