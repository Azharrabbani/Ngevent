export function formatEventTime(unix: number): string {
    const date = new Date(unix * 1000);
    return date.toLocaleTimeString("en-US", {
        hour: "2-digit",
        minute: "2-digit",
        hour12: false,
    });
}