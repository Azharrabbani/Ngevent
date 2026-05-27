import { useMutation, useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";
import { createCategoryApi } from "../api/categoryApi";
import { categoriesKeys } from "../../../utils/cacheKey";

export const useCreateCategory = (onSuccess?: () => void) => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: (name: string) => createCategoryApi(name),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: categoriesKeys.lists() });
            toast.success("Category created successfully!");
            onSuccess?.();
        },
        onError: (error: any) => {
            const message =
                error?.response?.data?.message ||
                error?.response?.data?.data ||
                "Failed to create category";
            toast.error(message);
        },
    });
};