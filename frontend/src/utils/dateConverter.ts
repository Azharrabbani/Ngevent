export const convertUnix = (unix: number): string => {
    let date: Date | string = new Date(unix * 1000);

    const day = date.getDate();
    const month = date.toLocaleString('en-US', { month: 'short' });
    const year = date.getFullYear();

    date = `${month} ${day} ${year}`;

    return date;
}

export const converDate = (date: Date | null): number => {
    if (!date) return 0;

    const utc = Date.UTC(
        date.getFullYear(),
        date.getMonth(),
        date.getDate()
    );

    return Math.floor(utc / 1000);
}

export const toDateString = (date: Date, month: "long" | "short" | "2-digit" | "numeric"): string => {
    const dateStr = date.toLocaleDateString("en-US", {
        timeZone: "Asia/Jakarta",
        weekday: "short",
        year: "numeric",
        month: month,
        day: "numeric",
    });

    return dateStr
}