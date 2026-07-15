// components/map/eventsMapContainer.tsx

import { useUserLocation } from "../../../hooks/useUserLocation";
import { SpinnerIcon } from "../../../../../components/icon";
import { useGetEventsForMap } from "../../../hooks/useGetEventForMap";
import EventsMapView from "./eventMapView";

interface Props {
    category?: number[];
    location?: string;
    search?: string;
    dateFilters?: Record<string, unknown>;
}

export default function EventsMapContainer({ category, location, search, dateFilters }: Props) {
    const { lat, lon } = useUserLocation();

    const { data, isLoading } = useGetEventsForMap({
        category,
        location,
        search,
        ...dateFilters,
        limit: 100,
    });

    if (isLoading) {
        return (
            <div className="w-full h-[70vh] flex items-center justify-center bg-white rounded-2xl border border-slate-200">
                <SpinnerIcon className="w-8 h-8 animate-spin text-blue-500" />
            </div>
        );
    }

    return (
        <div className="w-full h-[70vh]">
            <EventsMapView events={data?.rows ?? []} userLat={lat} userLon={lon} />
        </div>
    );
}