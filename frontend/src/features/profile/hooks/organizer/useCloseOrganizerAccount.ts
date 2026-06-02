import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import toast from "react-hot-toast";
import { CloseOrganizerAccountApi } from "../../api/profileApi";
import {
    organizerKeys,
    organizerUpdateKeys,
    eventsKeys,
    updateEventKeys,
    userKeys,
} from "../../../../utils/cacheKey";

export const useCloseOrganizerAccount = () => {
    const queryClient = useQueryClient();
    const navigate = useNavigate();

    return useMutation({
        mutationFn: CloseOrganizerAccountApi,
        onError: (err: any) => {
            toast.error(
                err?.response?.data?.error ||
                "Failed to close account. Please try again."
            );
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: userKeys.lists(), exact: false });
            queryClient.invalidateQueries({ queryKey: organizerKeys.all, exact: false });
            queryClient.invalidateQueries({ queryKey: organizerUpdateKeys.all, exact: false });
            queryClient.invalidateQueries({ queryKey: eventsKeys.all, exact: false });
            queryClient.invalidateQueries({ queryKey: updateEventKeys.all, exact: false });

            queryClient.clear();

            toast.success("Your account has been closed.");

            navigate("/login");
        },
    });
};