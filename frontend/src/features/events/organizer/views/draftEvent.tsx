import { useEffect, useState } from "react";
import { useGetCurrentOrganizerProfile } from "../../../profile/hooks/organizer/useGetCurrentOrganizerProfile";
import { useGetEvents } from "../hooks/useGetEvents";
import { converDate } from "../../../../utils/dateConverter";
import { defaultPagination } from "../../../../utils/pagination";
import { useListCategories } from "../../../categories/hooks/useListCategories";
import Sidebar from "../components/sidebar";
import Header from "../components/header";
import EventsContent from "../components/eventsContent";
import Pagination from "../../../../components/pagination";

export default function DraftEvents() {
    const [currentPage, setCurrentPage] = useState(1);
    const [location, setLocation] = useState<string | undefined>(undefined);
    const [event, setEvent] = useState<string | undefined>(undefined);
    const [selectedCategories, setSelectedCategories] = useState<number[]>([]);
    const [date, setDate] = useState<Date | null>(null);

    const organizer = useGetCurrentOrganizerProfile();

    const { data, isLoading } = useGetEvents({
        title: event,
        location: location,
        category: selectedCategories.length ? selectedCategories : undefined,
        status: "draft",
        start_time: date ? converDate(date) : undefined,
        pagination: defaultPagination(currentPage),
    });

    const { data: categoriesData, isLoading: categoriesLoading } = useListCategories();

    const totalPage = data?.total_pages || 1;


    const handleSearch = () => {
        setLocation(undefined);
        setEvent(undefined);
        setSelectedCategories([]);
        setDate(null);
        setCurrentPage(1);
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
                    date={date}
                    setDate={setDate}
                    setEvent={setEvent}
                    organizerName={organizer?.data?.name}
                    onSearch={handleSearch}
                    toggleStatus={false}
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