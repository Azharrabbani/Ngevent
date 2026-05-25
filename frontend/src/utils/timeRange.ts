export const timeRange = (startTime: Date, endTime: Date): string => {
    const options = { timeZone: "Asia/Jakarta" }

    return `${startTime.toLocaleTimeString("id-ID", { ...options, hour: "2-digit", minute: "2-digit" })} - ${endTime.toLocaleTimeString("id-ID", { ...options, hour: "2-digit", minute: "2-digit" })}`
};

export const formatTimeAgo = (date: Date | string) => {
    const now = new Date();
    const targetDate = new Date(date);

    const seconds = Math.floor((now.getTime() - targetDate.getTime()) / 1000);

    const intervals = [
        { label: "year", seconds: 31536000 },
        { label: "month", seconds: 2592000 },
        { label: "day", seconds: 86400 },
        { label: "hour", seconds: 3600 },
        { label: "minute", seconds: 60 },
        { label: "second", seconds: 1 },
    ];

    for (const interval of intervals) {
        const count = Math.floor(seconds / interval.seconds);

        if (count >= 1) {
            return count === 1
                ? `a ${interval.label} ago`
                : `${count} ${interval.label}s ago`;
        }
    }

    return "just now";
};