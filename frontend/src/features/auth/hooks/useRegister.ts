import type { RegisterRequest } from "../types/authRequest";
import { registerApi } from "../api/authApi";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { userKeys } from "../../../utils/cacheKey";
import toast from "react-hot-toast";

export const useRegister = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (payload: RegisterRequest) =>
            registerApi(payload),

        onSuccess: () => {
            toast.success("Admin registered successfully");

            queryClient.invalidateQueries({
                queryKey: userKeys.all,
                exact: false,
            });
        },

        onError: (err: any) => {
            const validationError =
                err?.response?.data?.errors;

            if (!validationError) {
                toast.error(
                    err?.response?.data?.error ||
                    "Failed register admin"
                );
            }
        },
    });
};