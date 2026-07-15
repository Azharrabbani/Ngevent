import type { RegisterRequest } from "../types/authRequest";
import { registerApi } from "../api/authApi";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { userKeys } from "../../../utils/cacheKey";
import toast from "react-hot-toast";
import { useState } from "react";

export const useRegister = () => {
    const queryClient = useQueryClient();

    const [errors, setErrors] = useState<Record<string, string>>({});
    const [generalError, setGeneralError] = useState<string>("");

    const mutation = useMutation({
        mutationFn: (payload: RegisterRequest) =>
            registerApi(payload),

        onMutate: () => {
            setErrors({});
            setGeneralError("");
        },

        onSuccess: (data, variables) => {
            localStorage.setItem(
                "verification_email",
                variables.email
            );

            toast.success(
                data.message
            );

            queryClient.invalidateQueries({
                queryKey: userKeys.all,
                exact: false,
            });
        },

        onError: (err: any) => {
            const validationErrors =
                err?.response?.data?.error;

            if (Array.isArray(validationErrors)) {
                const formattedErrors: Record<
                    string,
                    string
                > = {};

                validationErrors.forEach(
                    (e: any) => {
                        formattedErrors[e.field] =
                            e.message;
                    }
                );

                setErrors(formattedErrors);
                return;
            }

            setGeneralError(err?.response?.data?.error || "Failed register");

        },
    });

    return {
        ...mutation,
        generalError,
        errors,
    };
};