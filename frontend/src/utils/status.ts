export const getStatusColor = (status: string): string => {
    switch (status.toLowerCase()) {
        case "active":
            return "text-green-600"

        case "pending":
            return "text-yellow-600"

        case "rejected":
            return "text-red-600"

        case "done":
            return "text-blue-600"

        case "cancelled":
            return "text-gray-600"

        case "draft":
            return "text-purple-600"

        default:
            return "text-gray-500"
    }
}