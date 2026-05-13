import { logoutApi } from "../api/authApi"
import { useMutation, useQueryClient } from "@tanstack/react-query"

export const useLogout = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: () => logoutApi(),
        onSuccess: () => {
            queryClient.clear();
        },
    })
}