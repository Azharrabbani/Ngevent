import { useState, useCallback } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { FiArrowLeft } from "react-icons/fi";
import { useUserLocation } from "../../hooks/useUserLocation";
import { useGetEventBySlug } from "../../hooks/useGetEventBySlug";
import EventViewSkeleton from "../components/skeleton/eventViewSkeleton";
import LocationPermissionBanner from "../../../../components/locationPermission";
import EventBanner from "../components/event/eventBanner";
import EventAbout from "../components/event/eventAbout";
import LocationCard from "../components/cards/locationCard";
import EventOwnerCard from "../components/event_owner/eventOwnerCard";
import ShareCard from "../components/cards/shareCard";
import EventInfo from "../components/event/eventInfo";
import { ShockIcon } from "../../../../components/icon";
import EventMap from "../components/map/eventMap";

export default function EventViewPublicPage() {
    const { slug } = useParams<{ slug: string }>();
    const navigate = useNavigate();

    const { lat, lon, loading: locationLoading, denied } = useUserLocation();
    const [locationRequested, setLocationRequested] = useState(false);

    const handleRequestLocation = useCallback(() => {
        setLocationRequested(true);
        navigator.geolocation?.getCurrentPosition(
            () => window.location.reload(),
            () => setLocationRequested(false)
        );
    }, []);

    const {
        data: event,
        isLoading,
        isError,
    } = useGetEventBySlug(slug ?? "", {
        lat: lat,
        lon: lon,
    });

    if (isLoading || locationLoading) {
        return (
            <div className="min-h-screen bg-slate-50">
                <Header onBack={() => navigate(-1)} />
                <main className="max-w-6xl mx-auto px-4 lg:px-6 py-8">
                    <EventViewSkeleton />
                </main>
            </div>
        );
    }

    if (isError || !event) {
        return (
            <div className="min-h-screen bg-slate-50 flex flex-col">
                <Header onBack={() => navigate(-1)} />
                <div className="flex-1 flex flex-col items-center justify-center gap-3 text-center px-4">
                    <ShockIcon className="text-5xl text-indigo-500" />
                    <h2 className="text-lg font-semibold text-slate-800">Event not found</h2>
                    <p className="text-sm text-slate-500">
                        This event may have been removed or the link is invalid.
                    </p>
                    <button
                        onClick={() => navigate("/")}
                        className="mt-2 text-sm text-blue-600 font-medium hover:underline"
                    >
                        Back to events
                    </button>
                </div>
            </div>
        );
    }

    const showDistancePath =
        typeof event.distance === "string" &&
        Array.isArray(event.path) &&
        event.path.length > 0;

    return (
        <div className="min-h-screen bg-slate-50">
            <Header onBack={() => navigate(-1)} />

            <main className="max-w-6xl mx-auto px-4 lg:px-6 py-8 space-y-6">
                {denied && !locationRequested && (
                    <LocationPermissionBanner onRequestLocation={handleRequestLocation} />
                )}

                <div className="grid grid-cols-1 lg:grid-cols-[1fr_360px] gap-6 items-start">
                    <div className="space-y-5">
                        <EventBanner
                            banner={event.event.banner}
                            name={event.event.name}
                            categories={event.event.categories}
                        />

                        <EventAbout
                            description={event.event.description}
                        />

                        <LocationCard
                            address={event.event_address.address}
                            detailAddress={event.event_address.detail_address}
                            coordinates={event.event_address.coordinates}
                            path={showDistancePath ? event.path : undefined}
                            MapComponent={EventMap}
                        />
                    </div>

                    <div className="space-y-4 lg:sticky lg:top-6">
                        <EventInfo
                            name={event.event.name}
                            startTime={event.start_time}
                            endTime={event.end_time}
                            detailAddress={event.event_address.detail_address}
                            city={event.event_address.city}
                            distance={showDistancePath ? event.distance : undefined}
                        />

                        <EventOwnerCard
                            name={event.eo_profile.name}
                            photoProfile={event.eo_profile.photo_profile}
                            isVerified={event.eo_profile.is_verified}
                            onViewProfile={() =>
                                navigate(`/organizer/${event.eo_profile.id}`)
                            }
                        />

                        <ShareCard eventName={event.event.name} />
                    </div>
                </div>
            </main>
        </div>
    );
}

function Header({ onBack }: { onBack: () => void }) {
    return (
        <header className="bg-white border-b border-slate-200 sticky top-0 z-10">
            <div className="max-w-6xl mx-auto px-4 lg:px-6 h-14 flex items-center gap-4">
                <button
                    onClick={onBack}
                    className="flex items-center gap-1.5 text-sm text-slate-600 hover:text-slate-900 transition-colors"
                >
                    <FiArrowLeft className="w-4 h-4" />
                    Back
                </button>
            </div>
        </header>
    );
}