import EventFormContainer from "../components/form/eventFormContainer";
import EventForm from "../components/form/eventForm";
import { useDebounce } from "use-debounce";
import { useListCategories } from "../../../categories/hooks/useListCategories";
import { useState } from "react";
import { useSearchLocation } from "../../hooks/useSearchLocation";
import type { locationResp } from "../../types/locationResponse";

export default function CreateEvent() {
    const [selectedCategories, setSelectedCategories] = useState<number[]>([]);
    const [searchQuery, setSearchQuery] = useState("");
    const [selectedLocation, setSelectedLocation] = useState<locationResp | undefined>(undefined);

    const [debouncedQuery] = useDebounce(searchQuery, 500)
    
    const { data: categoriesData, isLoading: categoriesLoading } = useListCategories();
   
    const { data: locations, isLoading: locationLoading } = useSearchLocation(
        { query: debouncedQuery },
        {
            enabled: !!debouncedQuery
        }
    );
    
    
    return (
        <EventFormContainer>
            <EventForm
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
    )
}