import { useEffect } from "react";
import EventList from "../components/events/eventList";
import AdminSidebar from "../components/sideBar";
import { useGetEvents } from "../../events/hooks/useGetEvents";
import { useGetAllUpdatedEvents } from "../../events/hooks/useGetAllUpdatedEvent";
import { mapUpdateToEventResponse } from "../../events/utils/mapUpdateToEvent";
import type { EventsResponse } from "../../events/types/eventResponse";
import type { PaginatedData } from "../../../types/apiResponse";
import { buildEventDateFilters } from "../../../utils/dateFilter";
import { UseAdminEventFilters } from "../hooks/useAdminEventFilters";

export default function PendingEventList() {
    const {
        currentPage,
        setCurrentPage,

        search,
        setSearch,


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
    } = UseAdminEventFilters();

    const dateFilters = buildEventDateFilters(dateFilterType, selectedDate, selectedMonth, selectedYear);

    const isShowingUpdates = getUpdate === true;

    const { data: eventsData, isLoading: eventsLoading } = useGetEvents(
        {
            status: "pending",
            search,
            date: dateFilter,
            ...dateFilters,
            pagination: { limit: 6, page: currentPage },
        },
        !isShowingUpdates
    );

    const { data: updatesData, isLoading: updatesLoading } = useGetAllUpdatedEvents(
        {
            status: "pending",
            search,
            date: dateFilter,
            pagination: { limit: 6, page: currentPage },
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