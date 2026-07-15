import { useEffect, useState } from "react";
import { FiAlertCircle, FiNavigation } from "react-icons/fi";
import toast from "react-hot-toast";
import type { PathPoint } from "../../../types/publicEventResponse";
import { useGetEventRoute } from "../../../hooks/useGetEventRoute";
import type { UserLatLonRequest } from "../../../types/eventRequest";
import { MaximizeIcon, SpinnerIcon } from "../../../../../components/icon";
import MapModal from "../map/mapModal";


interface Props {
    eventId: string;
    address: string;
    detailAddress: string;
    coordinates: { lat: number; lon: number };
    path?: PathPoint[];
    distance?: string;
    userLocation: UserLatLonRequest;
    MapComponent: React.ComponentType<MapComponentProps>;
}

export interface MapComponentProps {
    center: [number, number];
    path?: PathPoint[];
}

export default function LocationCard({
    eventId,
    address,
    detailAddress,
    coordinates,
    path: initialPath,
    userLocation,
    MapComponent,
}: Props) {
    const [routeRequested, setRouteRequested] = useState(false);
    const [mapExpanded, setMapExpanded] = useState(false);

    const canFetchRoute = !!(userLocation.lat && userLocation.lon);

    const {
        data: routeData,
        isFetching: routeLoading,
        isError: routeError,
        errorUpdatedAt,
        refetch: refetchRoute,
    } = useGetEventRoute(eventId, userLocation, routeRequested && canFetchRoute);

    useEffect(() => {
        if (routeError) {
            toast.error("Failed to find route to event location, please try again later.");
        }
    }, [routeError, errorUpdatedAt]);

    const hasGeneratedRoute = Array.isArray(routeData?.path) && routeData.path.length > 0;
    const activePath = hasGeneratedRoute ? routeData.path : initialPath;
    const mapCenter: [number, number] = [coordinates.lat, coordinates.lon];

    const handleRouteButtonClick = () => {
        if (routeRequested) {
            refetchRoute();
        } else {
            setRouteRequested(true);
        }
    };

    return (
        <>
            <div className="bg-white rounded-2xl border border-slate-200 p-6 space-y-4">
                <h2 className="text-lg font-bold text-slate-900">Location</h2>

                <div
                    className={[
                        "relative rounded-xl overflow-hidden border border-slate-100 transition-all duration-500 ease-in-out",
                        hasGeneratedRoute ? "h-80" : "h-56",
                    ].join(" ")}
                >
                    <MapComponent center={mapCenter} path={activePath} />

                    <button
                        onClick={() => setMapExpanded(true)}
                        className="absolute top-2 right-2  flex items-center gap-1.5 px-2.5 py-1.5 bg-white/90 backdrop-blur-sm text-slate-700 text-xs font-medium rounded-lg border border-slate-200 shadow-sm hover:bg-white hover:text-slate-900 transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
                        aria-label="Expand map"
                    >
                        <MaximizeIcon className="w-3.5 h-3.5" />
                        <span>Expand</span>
                    </button>
                </div>

                <div>
                    <p className="text-sm font-semibold text-slate-800">{detailAddress}</p>
                    <p className="text-xs text-slate-500 mt-1 leading-relaxed">{address}</p>
                </div>

                {routeData?.distance && (
                    <p className="text-xs text-blue-600 font-medium">
                        {routeData.distance} from your location
                    </p>
                )}

                {canFetchRoute && (
                    <>
                        {!hasGeneratedRoute && (
                            <div className="space-y-2">
                                <button
                                    onClick={handleRouteButtonClick}
                                    disabled={routeLoading}
                                    className="w-full flex items-center justify-center gap-2 bg-slate-900 hover:bg-slate-700 disabled:bg-slate-400 text-white text-sm font-medium py-3 rounded-xl transition-colors"
                                >
                                    {routeLoading ? (
                                        <>
                                            <SpinnerIcon className="w-4 h-4 animate-spin" />
                                            Generating route...
                                        </>
                                    ) : routeError ? (
                                        <>
                                            <FiNavigation className="w-4 h-4" />
                                            Retry
                                        </>
                                    ) : (
                                        <>
                                            <FiNavigation className="w-4 h-4" />
                                            See Route
                                        </>
                                    )}
                                </button>

                                {routeError && !routeLoading && (
                                    <p className="flex items-center gap-1.5 text-xs text-red-600">
                                        <FiAlertCircle className="w-3.5 h-3.5 shrink-0" />
                                        Failed to find route to event location, please try again later.
                                    </p>
                                )}
                            </div>
                        )}

                        {hasGeneratedRoute && (
                            <button
                                onClick={() => {
                                    const dest = `${coordinates.lat},${coordinates.lon}`;
                                    window.open(
                                        `https://www.google.com/maps/dir/?api=1&destination=${dest}`,
                                        "_blank"
                                    );
                                }}
                                className="w-full flex items-center justify-center gap-2 border border-slate-200 hover:border-slate-300 text-slate-700 text-sm font-medium py-3 rounded-xl transition-colors"
                            >
                                <FiNavigation className="w-4 h-4" />
                                Open in Google Maps
                            </button>
                        )}
                    </>
                )}
            </div>

            <MapModal
                isOpen={mapExpanded}
                onClose={() => setMapExpanded(false)}
                center={mapCenter}
                path={activePath}
                MapComponent={MapComponent}
                title={detailAddress}
            />
        </>
    );
}