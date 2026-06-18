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

export function eventDateRange(
    startDateUnix: number,
    endDateUnix: number
) {
    const start = new Date(startDateUnix * 1000);
    const end = new Date(endDateUnix * 1000);

    const formatter = new Intl.DateTimeFormat(
        "en-GB",
        {
            day: "numeric",
            month: "short",
            year: "numeric",
        }
    );

    const shortFormatter = new Intl.DateTimeFormat(
        "en-GB",
        {
            day: "numeric",
            month: "short",
        }
    );


    const sameDay =
        start.getFullYear() === end.getFullYear() &&
        start.getMonth() === end.getMonth() &&
        start.getDate() === end.getDate();

    if (sameDay) {
        return toDateString(start, "short");
    }

    const sameYear =
        start.getFullYear() === end.getFullYear();

    if (sameYear) {
        return `${shortFormatter.format(start)} - ${formatter.format(end)}`;
    }

    return `${formatter.format(start)} - ${formatter.format(end)}`;
}