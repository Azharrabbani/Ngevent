import Sidebar from "../components/sidebar";
import Header from "../components/header";
import Pagination from "../../../../components/pagination";
import { useEffect } from "react";
import { useGetOrganizerEvents } from "../../hooks/useGetOrganizerEvents";
import { defaultPagination } from "../../../../utils/pagination";
import { useGetCurrentOrganizerProfile } from "../../../profile/hooks/organizer/useGetCurrentOrganizerProfile";
import { useListCategories } from "../../../categories/hooks/useListCategories";
import EventsContent from "../components/eventsContent";
import { buildEventDateFilters } from "../../../../utils/dateFilter";
import { useOrganizerEventFilters } from "../../hooks/useOrganizerEventFilters";

export default function CheckEvents() {
    const {
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
    } = useOrganizerEventFilters();

    const organizer = useGetCurrentOrganizerProfile();

    const dateFilters = buildEventDateFilters(
        dateFilterType,
        selectedDate,
        selectedMonth,
        selectedYear
    );

    const { data, isLoading } = useGetOrganizerEvents({
        title: event,
        location: location,
        category: selectedCategories.length ? selectedCategories : undefined,
        status: status,
        ...dateFilters,
        pagination: defaultPagination(currentPage),
    });

    const { data: categoriesData, isLoading: categoriesLoading } = useListCategories();

    const totalPage = data?.total_pages || 1;


    const handleSearch = () => {
        resetFilters
    };

    useEffect(() => {
        const delay = setTimeout(() => {
            setCurrentPage(1);
        }, 500);

        return () => clearTimeout(delay);
    }, [location, event]);

    return (
        <Sidebar photoProfile={organizer?.data?.photo_profile}>
            <>
                <Header
                    location={location}
                    setLocation={setLocation}
                    event={event}
                    categories={categoriesData}
                    categoriesLoading={categoriesLoading}
                    selectedCategories={selectedCategories}
                    setSelectedCategories={setSelectedCategories}
                    status={status}
                    dateFilterType={dateFilterType}
                    setDateFilterType={setDateFilterType}
                    selectedDate={selectedDate}
                    setSelectedDate={setSelectedDate}
                    selectedMonth={selectedMonth}
                    setSelectedMonth={setSelectedMonth}
                    selectedYear={selectedYear}
                    setSelectedYear={setSelectedYear}
                    setStatus={setStatus}
                    setEvent={setEvent}
                    organizerName={organizer?.data?.name}
                    onSearch={handleSearch}
                    toggleStatus={true}
                />

                <EventsContent
                    data={data}
                    loading={isLoading}
                />

                <Pagination
                    currentPage={currentPage}
                    totalPage={totalPage}
                    onPrev={() => setCurrentPage((prev) => Math.max(prev - 1, 1))}
                    onNext={() => setCurrentPage((prev) => Math.min(prev + 1, totalPage))}
                    onCurrent={(page) => setCurrentPage(page)}
                />

            </>
        </Sidebar>
    )
}