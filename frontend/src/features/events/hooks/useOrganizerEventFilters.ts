import { useSearchParams } from "react-router-dom";
import { useState, useEffect } from "react";
import type { DateFilterType } from "../../../utils/dateFilter";

export function useOrganizerEventFilters() {
    const [searchParams, setSearchParams] =
        useSearchParams();

    const [currentPage, setCurrentPage] =
        useState(
            Number(searchParams.get("page")) || 1
        );

    const [location, setLocation] =
        useState<string | undefined>(
            searchParams.get("location") ||
            undefined
        );

    const [event, setEvent] =
        useState<string | undefined>(
            searchParams.get("event") ||
            undefined
        );

    const [
        selectedCategories,
        setSelectedCategories,
    ] = useState<number[]>(
        searchParams.get("category")
            ? searchParams
                .get("category")!
                .split(",")
                .map(Number)
            : []
    );

    const [status, setStatus] =
        useState<string | undefined>(
            searchParams.get("status") ||
            undefined
        );

    const [
        dateFilterType,
        setDateFilterType,
    ] = useState<DateFilterType>(
        (searchParams.get(
            "dateType"
        ) as DateFilterType) ||
        "all"
    );

    const [selectedDate, setSelectedDate] =
        useState<Date | null>(
            searchParams.get("date")
                ? new Date(
                    searchParams.get(
                        "date"
                    )!
                )
                : null
        );

    const [selectedMonth, setSelectedMonth] =
        useState<number | undefined>(
            searchParams.get("month")
                ? Number(
                    searchParams.get(
                        "month"
                    )
                )
                : undefined
        );

    const [selectedYear, setSelectedYear] =
        useState<number | undefined>(
            searchParams.get("year")
                ? Number(
                    searchParams.get(
                        "year"
                    )
                )
                : undefined
        );

    useEffect(() => {
        const params =
            new URLSearchParams();

        if (currentPage > 1)
            params.set(
                "page",
                String(currentPage)
            );

        if (location)
            params.set(
                "location",
                location
            );

        if (event)
            params.set("event", event);

        if (selectedCategories.length)
            params.set(
                "category",
                selectedCategories.join(
                    ","
                )
            );

        if (status)
            params.set(
                "status",
                status
            );

        if (
            dateFilterType !== "all"
        ) {
            params.set(
                "dateType",
                dateFilterType
            );
        }

        if (selectedDate)
            params.set(
                "date",
                selectedDate.toISOString()
            );

        if (selectedMonth)
            params.set(
                "month",
                String(selectedMonth)
            );

        if (selectedYear)
            params.set(
                "year",
                String(selectedYear)
            );

        setSearchParams(params, {
            replace: true,
        });
    }, [
        currentPage,
        location,
        event,
        selectedCategories,
        status,
        dateFilterType,
        selectedDate,
        selectedMonth,
        selectedYear,
        setSearchParams,
    ]);

    const resetFilters = () => {
        setLocation(undefined);

        setEvent(undefined);

        setSelectedCategories([]);

        setStatus(undefined);

        setDateFilterType("all");

        setSelectedDate(null);

        setSelectedMonth(
            undefined
        );

        setSelectedYear(
            undefined
        );

        setCurrentPage(1);
    };

    return {
        currentPage,
        setCurrentPage,

        location,
        setLocation,

        event,
        setEvent,

        selectedCategories,
        setSelectedCategories,

        status,
        setStatus,

        dateFilterType,
        setDateFilterType,

        selectedDate,
        setSelectedDate,

        selectedMonth,
        setSelectedMonth,

        selectedYear,
        setSelectedYear,

        resetFilters,
    };
}