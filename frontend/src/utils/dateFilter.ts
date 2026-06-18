import { converDate } from "./dateConverter";

export type DateFilterType =
    | "all"
    | "today"
    | "this_month"
    | "this_year"
    | "date"
    | "month"
    | "year";


export function buildEventDateFilters(
    dateFilterType: DateFilterType,
    selectedDate: Date | null,
    selectedMonth?: number,
    selectedYear?: number,
) {
    const now = new Date();

    switch (dateFilterType) {
        case "today":
            return {
                event_date: converDate(now),
            };

        case "this_month":
            return {
                month: now.getMonth() + 1,
                year: now.getFullYear(),
            };

        case "this_year":
            return {
                year: now.getFullYear(),
            };

        case "date":
            return {
                event_date: selectedDate
                    ? converDate(selectedDate)
                    : undefined,
            };

        case "month":
            return {
                month: selectedMonth,
                year: selectedYear,
            };

        case "year":
            return {
                year: selectedYear,
            };

        default:
            return {};
    }
}