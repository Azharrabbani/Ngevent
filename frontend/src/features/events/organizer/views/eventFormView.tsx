import EventFormContainer from "../components/form/eventFormContainer";
import EventForm from "../components/form/eventForm";
import { useDebounce } from "use-debounce";
import { useListCategories } from "../../../categories/hooks/useListCategories";
import { useState } from "react";
import { useSearchLocation } from "../../hooks/useSearchLocation";
import type { locationResp } from "../../types/locationResponse";
import { useParams } from "react-router-dom";
import { useGetEventByID } from "../../hooks/useGetEventByID";


export default function EventFormView() {
    const { id } = useParams<{ id: string }>();
    const isEditMode = !!id;

    const [selectedCategories, setSelectedCategories] = useState<number[]>([]);
    const [searchQuery, setSearchQuery] = useState("");
    const [selectedLocation, setSelectedLocation] = useState<locationResp | undefined>(undefined);

    const [debouncedQuery] = useDebounce(searchQuery, 500);

    const { data: categoriesData, isLoading: categoriesLoading } = useListCategories();

    const { data: locations, isLoading: locationLoading } = useSearchLocation(
        { query: debouncedQuery },
        {
            enabled: !!debouncedQuery,
        }
    );

    // Only fetch event when in edit mode
    const { data: eventData, isLoading: eventLoading } = useGetEventByID(id ?? "");

    // Render the form until the event data is loaded in edit mode
    if (isEditMode && eventLoading) {
        return (
            <EventFormContainer>
                <div className="flex items-center justify-center min-h-[400px]">
                    <p className="text-gray-500 text-lg">Loading event...</p>
                </div>
            </EventFormContainer>
        );
    }

    return (
        <EventFormContainer>
            <EventForm
                mode={isEditMode ? "edit" : "create"}
                eventData={isEditMode ? eventData : undefined}
                categories={categoriesData}
                categoriesLoading={categoriesLoading}
                selectedCategories={selectedCategories}
                setSelectedCategories={setSelectedCategories}
                searchQuery={searchQuery}
                setSearchQuery={setSearchQuery}
                locations={locations}
                locationLoading={locationLoading}
                selectedLocation={selectedLocation}
                setSelectedLocation={setSelectedLocation}
            />
        </EventFormContainer>
    );
}