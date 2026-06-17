import { useSearchParams } from "react-router-dom";
import { useEffect, useState } from "react";
import type { DateFilterType } from "../../../utils/dateFilter";

export function UseAdminEventFilters() {
    const [searchParams, setSearchParams] =
        useSearchParams();

    const [currentPage, setCurrentPage] =
        useState(
            Number(searchParams.get("page")) || 1
        );

    const [search, setSearch] =
        useState<string | undefined>(
            searchParams.get("search") ||
            undefined
        );

    const [sort, setSort] =
        useState<string | undefined>(
            searchParams.get("sort") ||
            "desc"
        );

    const [dateFilter, setDateFilter] =
        useState<string | undefined>(
            searchParams.get("date") ||
            undefined
        );

    const [getUpdate, setGetUpdate] =
        useState<boolean | undefined>(
            searchParams.get(
                "update"
            ) === "true"
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
            searchParams.get(
                "selectedDate"
            )
                ? new Date(
                    searchParams.get(
                        "selectedDate"
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

        if (search)
            params.set(
                "search",
                search
            );

        if (sort)
            params.set(
                "sort",
                sort
            );

        if (dateFilter)
            params.set(
                "date",
                dateFilter
            );

        if (getUpdate)
            params.set(
                "update",
                "true"
            );

        if (
            dateFilterType !== "all"
        )
            params.set(
                "dateType",
                dateFilterType
            );

        if (selectedDate)
            params.set(
                "selectedDate",
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
        search,
        sort,
        dateFilter,
        getUpdate,
        dateFilterType,
        selectedDate,
        selectedMonth,
        selectedYear,
        setSearchParams,
    ]);

    useEffect(() => {
        const delay = setTimeout(() => {
            setCurrentPage(1);
        }, 500);

        return () =>
            clearTimeout(delay);
    }, [search, getUpdate]);

    return {
        currentPage,
        setCurrentPage,

        search,
        setSearch,

        sort,
        setSort,

        dateFilter,
        setDateFilter,

        getUpdate,
        setGetUpdate,

        dateFilterType,
        setDateFilterType,

        selectedDate,
        setSelectedDate,

        selectedMonth,
        setSelectedMonth,

        selectedYear,
        setSelectedYear,
    };
}