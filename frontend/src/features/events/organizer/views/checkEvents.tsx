import Sidebar from "../components/sidebar";
import Header from "../components/header";
import Pagination from "../../../../components/pagination";
import { useEffect, useState } from "react";
import { useGetEvents } from "../hooks/useGetEvents";
import { defaultPagination } from "../../../../utils/pagination";
import { useGetCurrentOrganizerProfile } from "../../../profile/hooks/organizer/useGetCurrentOrganizerProfile";
import { converDate } from "../../../../utils/dateConverter";
import { useListCategories } from "../../../categories/hooks/useListCategories";
import EventsContent from "../components/eventsContent";

export default function CheckEvents() {
    const [currentPage, setCurrentPage] = useState(1);
    const [location, setLocation] = useState<string | undefined>(undefined);
    const [event, setEvent] = useState<string | undefined>(undefined);
    const [selectedCategories, setSelectedCategories] = useState<number[]>([]);
    const [status, setStatus] = useState<string | undefined>(undefined);
    const [date, setDate] = useState<Date | null>(null);

    const organizer = useGetCurrentOrganizerProfile();

    const { data, isLoading } = useGetEvents({    
        title: event,
        location: location,
        category: selectedCategories.length ? selectedCategories : undefined,
        status: status,
        date: date ? converDate(date) : undefined,
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
                    status={status}
                    date={date}
                    setDate={setDate}
                    setStatus={setStatus}
                    setEvent={setEvent}
                    organizerName={organizer?.data?.name}
                    onSearch={handleSearch}
                    toggleStatus= {true}
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