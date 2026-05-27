import { useMutation, useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";
import { updateCategoryApi } from "../api/categoryApi";
import { categoriesKeys } from "../../../utils/cacheKey";

export const useUpdateCategory = () => {
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: ({ id, name }: { id: string | number; name: string }) =>
            updateCategoryApi({ id, name }),
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: categoriesKeys.lists() });
            toast.success("Category updated successfully!");
        },
        onError: (error: any) => {
            const message =
                error?.response?.data?.message ||
                error?.response?.data?.data ||
                "Failed to update category";
            toast.error(message);
        },
    });
};