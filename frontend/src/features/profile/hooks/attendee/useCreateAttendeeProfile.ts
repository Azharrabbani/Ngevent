import { createAttendeeProfileApi } from "../../api/profileApi"
import toast from "react-hot-toast";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { attendeeKeys } from "../../utils/cacheKey";

export const useCreateAttendeeProfile = () => {
    const queryClient = useQueryClient();
    
     return useMutation({
        mutationFn: createAttendeeProfileApi,

        onError: (err: any) => {
            const validationError = err?.response?.data?.error;

            if (!Array.isArray(validationError)) {
                toast.error(err?.response?.data?.error || "Failed create profile");
            }
        },

        onSuccess: () => {
            toast.success("Profile created successfully");

            queryClient.invalidateQueries({ queryKey: attendeeKeys.lists() });
            queryClient.invalidateQueries({ queryKey: attendeeKeys.details() });
            queryClient.invalidateQueries({ queryKey: attendeeKeys.me() });
            queryClient.invalidateQueries({ queryKey: [""] });
        }
    });
}