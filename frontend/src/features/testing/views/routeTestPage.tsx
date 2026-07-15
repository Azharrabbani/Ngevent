import { useState } from "react";
import { SpinnerIcon, MaximizeIcon } from "../../../components/icon";
import { useRouteTest } from "../../events/hooks/useGetRouteTest";
import MapModal from "../../events/public/components/map/mapModal";
import EventMap from "../../events/public/components/map/eventMap";


const USERS = [
    { name: "User-Condet", lat: -6.2666011, lon: 106.8624308 },
    { name: "User-PasarMinggu", lat: -6.3144585, lon: 106.8268601 },
    { name: "User-PondokLabu", lat: -6.305911, lon: 106.7994565 },
    { name: "User-TanjungPriok", lat: -6.1330885, lon: 106.8695885 },
    { name: "User-Kemayoran", lat: -6.1606132, lon: 106.8422236 },
    { name: "User-PasarRebo", lat: -6.3313338, lon: 106.845284 },
    { name: "User-Pluit", lat: -6.1109418, lon: 106.7781484 },
    { name: "User-Semanggi", lat: -6.226084, lon: 106.8202041 },
    { name: "User-Cakung", lat: -6.1727261, lon: 106.9354858 },
    { name: "User-Ragunan", lat: -6.3175008, lon: 106.8269633 },
];

const EVENTS = [
    { name: "Event-GrandIndonesia", lat: -6.1957601, lon: 106.8214547 },
    { name: "Event-Ancol", lat: -6.125215, lon: 106.8362474 },
    { name: "Event-TMI", lat: -6.3026341, lon: 106.8952962 },
    { name: "Event-JIS", lat: -6.1250747, lon: 106.8609623 },
    { name: "Event-GBK", lat: -6.2186492, lon: 106.8036258 },
    { name: "Event-KotaKasablanka", lat: -6.2232551, lon: 106.8426972 },
    { name: "Event-TebetEcoPark", lat: -6.2403737, lon: 106.8524328 },
    { name: "Event-SummareconMallBekasi", lat: -6.2260106, lon: 107.0010614 },
    { name: "Event-BSDTangerang", lat: -6.300681, lon: 106.6365721 },
    { name: "Event-MargoCity", lat: -6.3728858, lon: 106.8350994 },
    
];


export default function RouteTestPage() {
    const [userIdx, setUserIdx] = useState(0);
    const [eventIdx, setEventIdx] = useState(0);
    const [mapExpanded, setMapExpanded] = useState(false);

    const user = USERS[userIdx];
    const event = EVENTS[eventIdx];

    const { data, isFetching, refetch, isError } = useRouteTest({
        from_lat: user.lat,
        from_lon: user.lon,
        to_lat: event.lat,
        to_lon: event.lon,
        to_name: event.name,
    });

    const hasRoute = Array.isArray(data?.path) && data.path.length > 0;
    const mapCenter: [number, number] = [event.lat, event.lon];

    const googleMapsUrl = data
        ? `https://www.google.com/maps/dir/?api=1&origin=${user.lat},${user.lon}&destination=${event.lat},${event.lon}`
        : undefined;

    return (
        <div className="max-w-2xl mx-auto p-6 space-y-6">
            <h1 className="text-2xl font-bold text-slate-900">Route Test — Jabodetabek Dataset</h1>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <div>
                    <label className="block text-sm font-medium text-slate-700 mb-1">Starting Point (User)</label>
                    <select
                        value={userIdx}
                        onChange={(e) => setUserIdx(Number(e.target.value))}
                        className="w-full border border-slate-300 rounded-lg px-3 py-2 text-sm"
                    >
                        {USERS.map((u, i) => (
                            <option key={u.name} value={i}>
                                {u.name} ({u.lat.toFixed(4)}, {u.lon.toFixed(4)})
                            </option>
                        ))}
                    </select>
                </div>

                <div>
                    <label className="block text-sm font-medium text-slate-700 mb-1">Ending Point (Event)</label>
                    <select
                        value={eventIdx}
                        onChange={(e) => setEventIdx(Number(e.target.value))}
                        className="w-full border border-slate-300 rounded-lg px-3 py-2 text-sm"
                    >
                        {EVENTS.map((ev, i) => (
                            <option key={ev.name} value={i}>
                                {ev.name} ({ev.lat.toFixed(4)}, {ev.lon.toFixed(4)})
                            </option>
                        ))}
                    </select>
                </div>
            </div>

            <button
                onClick={() => refetch()}
                disabled={isFetching}
                className="w-full flex items-center justify-center gap-2 bg-slate-900 hover:bg-slate-700 disabled:bg-slate-400 text-white text-sm font-medium py-3 rounded-xl transition-colors"
            >
                {isFetching ? (
                    <>
                        <SpinnerIcon className="w-4 h-4 animate-spin" />
                        Calculating route...
                    </>
                ) : (
                    "Calculate Route"
                )}
            </button>

            {isError && (
                <p className="text-sm text-red-600">Failed to calculate route. Try again.</p>
            )}

            {data && (
                <div className="border border-slate-200 rounded-2xl p-5 space-y-3 bg-white">
                    <h2 className="font-semibold text-slate-900">Calculation Results</h2>

                    <div className="grid grid-cols-2 gap-3 text-sm">
                        <div className="bg-slate-50 rounded-lg p-3">
                            <p className="text-xs text-slate-500">Route Distance</p>
                            <p className="font-semibold text-slate-900">{data.distance}</p>
                        </div>
                        <div className="bg-slate-50 rounded-lg p-3">
                            <p className="text-xs text-slate-500">Dijkstra Cost</p>
                            <p className="font-semibold text-slate-900">{data.dijkstra_cost.toFixed(4)}</p>
                        </div>
                        <div className="bg-slate-50 rounded-lg p-3">
                            <p className="text-xs text-slate-500">Dijkstra Execution Time</p>
                            <p className="font-semibold text-slate-900">{data.dijkstra_time_ms.toFixed(2)} ms</p>
                        </div>
                        <div className="bg-slate-50 rounded-lg p-3">
                            <p className="text-xs text-slate-500">Waypoints</p>
                            <p className="font-semibold text-slate-900">{data.path.length} titik</p>
                        </div>
                    </div>

                    {hasRoute && (
                        <div className="relative h-80 rounded-xl overflow-hidden border border-slate-100">
                            <EventMap center={mapCenter} path={data.path} />

                            <button
                                onClick={() => setMapExpanded(true)}
                                className="absolute top-2 right-2 flex items-center gap-1.5 px-2.5 py-1.5 bg-white/90 backdrop-blur-sm text-slate-700 text-xs font-medium rounded-lg border border-slate-200 shadow-sm hover:bg-white hover:text-slate-900 transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
                                aria-label="Expand map"
                            >
                                <MaximizeIcon className="w-3.5 h-3.5" />
                                <span>Expand</span>
                            </button>
                        </div>
                    )}

                    <a href={googleMapsUrl}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="block w-full text-center bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium py-3 rounded-xl transition-colors"
                    >
                        Open in Google Maps ({user.name} → {event.name})
                    </a>
                </div>
            )
            }

            <MapModal
                isOpen={mapExpanded}
                onClose={() => setMapExpanded(false)}
                center={mapCenter}
                path={data?.path}
                MapComponent={EventMap}
                title={`${user.name} → ${event.name}`}
            />
        </div >
    );
}