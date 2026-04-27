import { useMutation, useQueryClient } from "@tanstack/react-query"
import type { rejectOrganizerReq } from "../../types/profileRequest";
import { RejectOrganizerApi } from "../../api/profileApi";
import toast from "react-hot-toast";
import { organizerKeys } from "../../utils/cacheKey";

export const useRejectOrganizer = (id: string) => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: ({ id, payload }: { id: string; payload: rejectOrganizerReq }) => {
            return RejectOrganizerApi(id, payload);
        },
        onError: (err: any) => {
            toast.error(err?.response?.data?.error || "Rejection failed");
        },
        onSuccess: (success) => {
            toast.success(success.data);

            queryClient.invalidateQueries({
                queryKey: organizerKeys.me(),
            });

             queryClient.invalidateQueries({
                queryKey: organizerKeys.detail(id),
            });

            queryClient.invalidateQueries({
                queryKey: organizerKeys.lists(),
                exact: false,
            });
        }
    })
}