import { useMutation, useQueryClient } from "@tanstack/react-query"
import { ApproveOrganizerApi } from "../../api/profileApi";
import toast from "react-hot-toast";
import { organizerKeys } from "../../utils/cacheKey";

export const useApproveOrganizer = (id: string) => {
    const queryClient = useQueryClient()    ;

    return useMutation({
        mutationFn: ApproveOrganizerApi,
        onError: (err: any) => {
            toast.error(err?.response?.data?.error || "Approval failed");
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