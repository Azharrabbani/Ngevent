import AdminSidebar from "../components/sideBar";
import EventList from "../components/events/eventList";
import { useEffect, useState } from "react";
import { useGetEvents } from "../../events/hooks/useGetEvents";
import { buildEventDateFilters, type DateFilterType } from "../../../utils/dateFilter";

export default function ActiveEventList() {
    const [currentPage, setCurrentPage] = useState(1);
    const [search, setSearch] = useState<string | undefined>(undefined);
    const [sort, setSort] = useState<string | undefined>("desc");
    const [filter, setFilter] = useState<string | undefined>(undefined);
    const [dateFilterType, setDateFilterType] = useState<DateFilterType>("all");
    const [selectedDate, setSelectedDate] = useState<Date | null>(null);
    const [selectedMonth, setSelectedMonth] = useState<number | undefined>();
    const [selectedYear, setSelectedYear] = useState<number | undefined>();

    const dateFilters = buildEventDateFilters(dateFilterType, selectedDate, selectedMonth, selectedYear);


    const { data, isLoading } = useGetEvents({
        status: "active",
        search: search,
        sort: sort,
        date: filter,
        ...dateFilters,
        pagination: {
            limit: 4,
            page: currentPage,
        }
    })

    const totalPage = data?.total_pages ?? 1;

    useEffect(() => {
        const delay = setTimeout(() => {
            setCurrentPage(1);
        }, 500);

        return () => clearTimeout(delay);
    }, [search]);

    return (
        <AdminSidebar>
            <>
                <EventList
                    status="active"
                    data={data}
                    isLoading={isLoading}
                    reviewEvent={false}
                    currentPage={currentPage}
                    totalPage={totalPage}
                    search={search}
                    setSearch={setSearch}
                    sort={sort}
                    setSort={setSort}
                    dateFilter={filter}
                    setDateFilter={setFilter}
                    dateFilterType={dateFilterType}
                    setDateFilterType={setDateFilterType}
                    selectedDate={selectedDate}
                    setSelectedDate={setSelectedDate}
                    selectedMonth={selectedMonth}
                    setSelectedMonth={setSelectedMonth}
                    selectedYear={selectedYear}
                    setSelectedYear={setSelectedYear}
                    setCurrentPage={setCurrentPage}
                />
            </>
        </AdminSidebar>
    )
}