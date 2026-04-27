import { useMutation, useQueryClient } from "@tanstack/react-query"
import { ValidateOrganizerUpdateApi } from "../../api/profileApi";
import type { validateOrganizerReq } from "../../types/profileRequest";
import { organizerKeys, organizerUpdateKeys } from "../../utils/cacheKey";
import toast from "react-hot-toast";

export const useValidateOrganizerUpdate = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: ({ id, payload }: { id: string; payload: validateOrganizerReq }) => {
            return ValidateOrganizerUpdateApi(id, payload);
        },
        onMutate: async ({ id, payload }) => {
                await queryClient.cancelQueries({
                queryKey: organizerKeys.detail(id)
            });

            const previousDetail = queryClient.getQueryData(
                organizerKeys.detail(id)
            );

            queryClient.setQueryData(
                organizerKeys.detail(id),
                (old: any) => {
                    if (!old) return old;

                    return {
                        ...old,
                        status: {
                            ...old.status,
                            status: payload.status,
                        },
                    };
                }
            );

            return { previousDetail, id };
        },
        onError: (err: any, _, context) => {
            toast.error(err?.response?.data?.error || "Validation failed");

            if (context?.previousDetail) {
                queryClient.setQueryData(
                    organizerKeys.detail(context.id),
                    context.previousDetail
                );
            }
        },
        onSuccess: (_, { id, payload }) => {
            if (payload.status === "approved") {
                toast.success("Organizer approved");
            } else {
                toast.success("Organizer rejected");
            }

            queryClient.invalidateQueries({
                queryKey: organizerKeys.detail(id),
            });

            queryClient.invalidateQueries({
                queryKey: organizerUpdateKeys.detail(id),
            });

            queryClient.invalidateQueries({
                queryKey: organizerKeys.lists(),
                exact: false,
                refetchType: "active"
            });
        },
    });
};