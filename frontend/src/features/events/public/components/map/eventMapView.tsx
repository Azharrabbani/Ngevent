// components/map/eventsMapView.tsx
import { useEffect, useMemo } from "react";
import {
    MapContainer,
    TileLayer,
    Marker,
    Popup,
    useMap,
} from "react-leaflet";
import L from "leaflet";
import "leaflet/dist/leaflet.css";
import markerIcon2x from "leaflet/dist/images/marker-icon-2x.png";
import markerIcon from "leaflet/dist/images/marker-icon.png";
import markerShadow from "leaflet/dist/images/marker-shadow.png";
import type { EventsResponse } from "../../../types/eventResponse";
import { eventDateRange } from "../../../../../utils/dateConverter";

delete (L.Icon.Default.prototype as any)._getIconUrl;
L.Icon.Default.mergeOptions({
    iconRetinaUrl: markerIcon2x,
    iconUrl: markerIcon,
    shadowUrl: markerShadow,
});

const userIcon = L.divIcon({
    className: "",
    html: `
        <div style="
            width: 18px; height: 18px;
            background: #2563eb;
            border: 3px solid white;
            border-radius: 50%;
            box-shadow: 0 0 0 2px rgba(37,99,235,0.4);
        "></div>
    `,
    iconSize: [18, 18],
    iconAnchor: [9, 9],
});

interface EventsMapViewProps {
    events: EventsResponse[];
    userLat?: number | null;
    userLon?: number | null;
}

function FitBounds({ positions }: { positions: [number, number][] }) {
    const map = useMap();

    useEffect(() => {
        if (positions.length === 0) return;

        if (positions.length === 1) {
            map.setView(positions[0], 14);
            return;
        }

        const bounds = L.latLngBounds(positions);
        map.fitBounds(bounds, { padding: [40, 40] });
    }, [positions, map]);

    return null;
}

export default function EventsMapView({ events, userLat, userLon }: EventsMapViewProps) {
    const eventPoints = useMemo(
        () =>
            events
                .filter((e) => e.event_address?.coordinates?.lat && e.event_address?.coordinates?.lon)
                .map((e) => ({
                    id: e.id,
                    name: e.event.name,
                    banner: e.event.banner,
                    city: e.event_address.city,
                    slug: e.event.slug,
                    dateLabel: eventDateRange(e.start_date, e.end_date),
                    position: [
                        e.event_address.coordinates.lat,
                        e.event_address.coordinates.lon,
                    ] as [number, number],
                })),
        [events]
    );

    const allPositions: [number, number][] = [
        ...eventPoints.map((e) => e.position),
        ...(userLat && userLon ? [[userLat, userLon] as [number, number]] : []),
    ];

    const defaultCenter: [number, number] = allPositions[0] || [-6.2, 106.816666]; // fallback Jakarta

    return (
        <MapContainer
            center={defaultCenter}
            zoom={12}
            scrollWheelZoom
            className="w-full h-full rounded-2xl"
        >
            <TileLayer
                attribution='&copy; OpenStreetMap contributors'
                url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            />

            <FitBounds positions={allPositions} />

            {eventPoints.map((event) => (
                <Marker key={event.id} position={event.position}>
                    <Popup>
                        <div className="w-48">
                            {event.banner && (
                                <img
                                    src={event.banner}
                                    alt={event.name}
                                    className="w-full h-24 object-cover rounded-md mb-2"
                                />
                            )}
                            <p className="font-semibold text-sm text-slate-900">{event.name}</p>
                            <p className="text-xs text-slate-500">{event.city}</p>
                            <p className="text-xs text-slate-500 mt-0.5">{event.dateLabel}</p>


                            <a href={`/events/${event.slug}`}
                                className="inline-block mt-2 text-xs font-medium text-blue-600 hover:underline"
                            >
                                Lihat Detail →
                            </a>
                        </div>
                    </Popup>
                </Marker>
            ))
            }

            {
                userLat && userLon && (
                    <Marker position={[userLat, userLon]} icon={userIcon}>
                        <Popup>Lokasi kamu</Popup>
                    </Marker>
                )
            }
        </MapContainer >
    );
}