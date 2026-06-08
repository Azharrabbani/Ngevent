import { useEffect, useState } from "react";
import EventList from "../components/events/eventList";
import AdminSidebar from "../components/sideBar";
import { useGetEvents } from "../../events/hooks/useGetEvents";
import { useGetAllUpdatedEvents } from "../../events/hooks/useGetAllUpdatedEvent";
import { mapUpdateToEventResponse } from "../../events/utils/mapUpdateToEvent";
import type { EventsResponse } from "../../events/types/eventResponse";
import type { PaginatedData } from "../../../types/apiResponse";
import { buildEventDateFilters, type DateFilterType } from "../../../utils/dateFilter";

export default function PendingEventList() {
    const [currentPage, setCurrentPage] = useState(1);
    const [search, setSearch] = useState<string | undefined>(undefined);
    const [sort, setSort] = useState<string | undefined>("desc");
    const [dateFilter, setDateFilter] = useState<string | undefined>(undefined);
    const [getUpdate, setGetUpdate] = useState<boolean | undefined>(undefined);
    const [dateFilterType, setDateFilterType] = useState<DateFilterType>("all");
    const [selectedDate, setSelectedDate] = useState<Date | null>(null);
    const [selectedMonth, setSelectedMonth] = useState<number | undefined>();
    const [selectedYear, setSelectedYear] = useState<number | undefined>();

    const dateFilters = buildEventDateFilters(dateFilterType, selectedDate, selectedMonth, selectedYear);

    const isShowingUpdates = getUpdate === true;

    const { data: eventsData, isLoading: eventsLoading } = useGetEvents(
        {
            status: "pending",
            search,
            sort,
            date: dateFilter,
            ...dateFilters,
            pagination: { limit: 4, page: currentPage },
        },
        !isShowingUpdates
    );

    const { data: updatesData, isLoading: updatesLoading } = useGetAllUpdatedEvents(
        {
            status: "pending",
            search,
            sort,
            date: dateFilter,
            pagination: { limit: 4, page: currentPage },
        },
        isShowingUpdates
    );

    const mappedUpdatesData: PaginatedData<EventsResponse> | undefined =
        isShowingUpdates && updatesData
            ? {
                ...updatesData,
                rows: updatesData.rows?.map(mapUpdateToEventResponse) ?? [],
            }
            : undefined;

    const data = isShowingUpdates ? mappedUpdatesData : eventsData;
    const isLoading = isShowingUpdates ? updatesLoading : eventsLoading;
    const totalPage = data?.total_pages ?? 1;

    useEffect(() => {
        const delay = setTimeout(() => setCurrentPage(1), 500);
        return () => clearTimeout(delay);
    }, [search, getUpdate]);

    return (
        <AdminSidebar>
            <>
                <EventList
                    status="pending"
                    data={data}
                    isLoading={isLoading}
                    reviewEvent={true}
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
                    getUpdate={getUpdate}
                    setGetupdate={setGetUpdate}
                    setCurrentPage={setCurrentPage}
                />
            </>
        </AdminSidebar>
    );
}