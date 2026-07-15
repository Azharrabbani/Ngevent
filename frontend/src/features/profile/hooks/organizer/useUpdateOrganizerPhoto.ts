import { useMutation, useQueryClient } from "@tanstack/react-query"

import toast from "react-hot-toast";
import { UpdateOrganizerPhotoApi } from "../../api/profileApi";
import { organizerKeys } from "../../../../utils/cacheKey";

export const useUpdateOrganizerPhoto = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: UpdateOrganizerPhotoApi,
        onMutate: async (newData) => {
            await queryClient.cancelQueries({ queryKey: organizerKeys.me() });

            const previous = queryClient.getQueryData(organizerKeys.me());

            const previewUrl = URL.createObjectURL(newData.photo);

            queryClient.setQueryData(organizerKeys.me(), (old: any) => ({
                ...old,
                photo_profile: previewUrl,
            }));

            return { previous, previewUrl };
        },
        onError: (err: any, _, context) => {
            toast.error(err?.response?.data?.error || "Update failed")
            queryClient.setQueryData(organizerKeys.me(), context?.previous);

            if (context?.previewUrl) {
                URL.revokeObjectURL(context.previewUrl);
            }
        },
        onSuccess: (_, __, context) => {
            if (context?.previewUrl) {
                URL.revokeObjectURL(context.previewUrl);
            }

            queryClient.invalidateQueries({ queryKey: organizerKeys.me() });
        },
    });
}