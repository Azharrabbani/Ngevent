import type { RegisterRequest } from "../types/authRequest";
import { registerAdminApi } from "../api/authApi";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { userKeys } from "../../../utils/cacheKey";
import toast from "react-hot-toast";

export const useRegisterAdmin = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (payload: RegisterRequest) =>
            registerAdminApi(payload),

        onSuccess: () => {
            toast.success("Admin registered successfully");

            queryClient.invalidateQueries({
                queryKey: userKeys.all,
                refetchType: "active"

            });
        },

        onError: (err: any) => {
            const validationError = err?.response?.data?.error;

            if (!Array.isArray(validationError)) {
                toast.error(err?.response?.data?.error || "Failed register admin");
            }
        },
    });
};