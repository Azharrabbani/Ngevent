import { useEffect } from "react";
import EventList from "../components/events/eventList";
import AdminSidebar from "../components/sideBar";
import { useGetEvents } from "../../events/hooks/useGetEvents";
import { buildEventDateFilters } from "../../../utils/dateFilter";
import { UseAdminEventFilters } from "../hooks/useAdminEventFilters";

export default function DoneEventList() {
    const {
        currentPage,
        setCurrentPage,

        search,
        setSearch,

        sort,
        setSort,

        dateFilter,
        setDateFilter,

        dateFilterType,
        setDateFilterType,

        selectedDate,
        setSelectedDate,

        selectedMonth,
        setSelectedMonth,

        selectedYear,
        setSelectedYear,
    } = UseAdminEventFilters();

    const dateFilters = buildEventDateFilters(dateFilterType, selectedDate, selectedMonth, selectedYear);


    const { data, isLoading } = useGetEvents({
        status: "done",
        search: search,
        sort: sort,
        date: dateFilter,
        ...dateFilters,
        pagination: {
            limit: 6,
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
                    status="done"
                    data={data}
                    isLoading={isLoading}
                    reviewEvent={false}
                    currentPage={currentPage}
                    totalPage={totalPage}
                    search={search}
                    setSearch={setSearch}
                    sort={sort}
                    setSort={setSort}
                    dateFilter={dateFilter}
                    setDateFilter={setDateFilter}
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