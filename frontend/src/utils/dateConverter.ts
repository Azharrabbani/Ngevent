export const convertUnix = (unix: number): string => {
    let date: Date | string = new Date(unix * 1000);

    const day = date.getDate();
    const month = date.toLocaleString('en-US', { month: 'short' });
    const year = date.getFullYear();

    date = `${month} ${day} ${year}`;
    
    return date;
}

export const converDate = (date: Date): number => {
    return Math.floor(date.getTime() / 1000)
}