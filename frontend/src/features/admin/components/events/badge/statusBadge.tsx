interface StatusBadgeProps {
    status: string;
    className?: string;
}

const statusConfig: Record<
    string,
    { label: string; bg: string; text: string }
> = {
    ACTIVE: {
        label: "ACTIVE",
        bg: "bg-green-100",
        text: "text-green-700",
    },
    DONE: {
        label: "DONE",
        bg: "bg-blue-100",
        text: "text-blue-700",
    },
    PENDING: {
        label: "PENDING",
        bg: "bg-amber-100",
        text: "text-amber-700",
    },
    REJECTED: {
        label: "REJECTED",
        bg: "bg-red-100",
        text: "text-red-600",
    },
    UPDATE_REQUEST: {
        label: "UPDATE REQUEST",
        bg: "bg-amber-500",
        text: "text-white",
    },
    NEW_REQUEST: {
        label: "NEW REQUEST",
        bg: "bg-blue-500",
        text: "text-white",
    },
};

export default function StatusBadge({ status, className = "" }: StatusBadgeProps) {
    const config = statusConfig[status];
    return (
        <span
            className={`inline-flex items-center px-2.5 py-1 rounded text-xs font-bold tracking-wide uppercase ${config.bg} ${config.text} ${className}`}
        >
            {config.label}
        </span>
    );
}

