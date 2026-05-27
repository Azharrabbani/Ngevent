export function FormatRelativeTime(unixTime: number): string {
    const now = new Date();
    const date = new Date(unixTime * 1000);

    const diff = now.getTime() - date.getTime();

    const seconds = Math.floor(diff / 1000);
    const minutes = Math.floor(diff / (1000 * 60));
    const hours = Math.floor(diff / (1000 * 60 * 60));
    const days = Math.floor(diff / (1000 * 60 * 60 * 24));
    const weeks = Math.floor(days / 7);

    const months =
        (now.getFullYear() - date.getFullYear()) * 12 +
        (now.getMonth() - date.getMonth());

    const years = now.getFullYear() - date.getFullYear();

    if (seconds < 10) {
        return "now";
    }

    if (minutes < 60) {
        return minutes === 1
            ? "1 minute ago"
            : `${minutes} minutes ago`;
    }

    if (hours < 24) {
        return hours === 1
            ? "1 hour ago"
            : `${hours} hours ago`;
    }

    if (days === 1) {
        return "yesterday";
    }

    if (days < 14) {
        return `${days} days ago`;
    }

    if (weeks < 5) {
        return weeks === 1
            ? "a week ago"
            : `${weeks} weeks ago`;
    }

    if (months < 12) {
        return months === 1
            ? "a month ago"
            : `${months} months ago`;
    }

    return years === 1
        ? "a year ago"
        : `${years} years ago`;
}