import { useEffect, useRef } from "react";
import { MapContainer, TileLayer, Marker, Popup, Polyline, useMap } from "react-leaflet";
import L from "leaflet";
import markerIcon2x from "leaflet/dist/images/marker-icon-2x.png";
import markerIcon from "leaflet/dist/images/marker-icon.png";
import markerShadow from "leaflet/dist/images/marker-shadow.png";
import type { MapComponentProps } from "../cards/locationCard";
import type { PathPoint } from "../../../types/publicEventResponse";

delete (L.Icon.Default.prototype as any)._getIconUrl;
L.Icon.Default.mergeOptions({
    iconUrl: markerIcon,
    iconRetinaUrl: markerIcon2x,
    shadowUrl: markerShadow,
});

const eventIcon = new L.Icon({
    iconUrl: markerIcon,
    iconRetinaUrl: markerIcon2x,
    shadowUrl: markerShadow,
    iconSize: [25, 41],
    iconAnchor: [12, 41],
    popupAnchor: [1, -34],
    shadowSize: [41, 41],
});

const userIcon = new L.DivIcon({
    className: "",
    html: `
        <div style="
            width: 16px;
            height: 16px;
            background: #3b82f6;
            border: 3px solid white;
            border-radius: 50%;
            box-shadow: 0 2px 6px rgba(59,130,246,0.5);
        "></div>
    `,
    iconSize: [16, 16],
    iconAnchor: [8, 8],
    popupAnchor: [0, -10],
});

function FitBounds({ points }: { points: [number, number][] }) {
    const map = useMap();
    const fittedRef = useRef(false);

    useEffect(() => {
        if (fittedRef.current || points.length < 2) return;
        const bounds = L.latLngBounds(points.map(([lat, lon]) => [lat, lon]));
        map.fitBounds(bounds, { padding: [40, 40] });
        fittedRef.current = true;
    }, [map, points]);

    return null;
}

function ChangeCenter({ center }: { center: [number, number] }) {
    const map = useMap();
    useEffect(() => {
        map.setView(center, map.getZoom());
    }, [map, center]);
    return null;
}

export default function EventMap({ center, path }: MapComponentProps) {
    const hasPath = Array.isArray(path) && path.length >= 2;
    const pathCoords: [number, number][] = hasPath
        ? (path as PathPoint[]).map((p) => [p.lat, p.lon])
        : [];

    const userPoint = hasPath ? (path as PathPoint[])[0] : null;
    const eventPoint = hasPath ? (path as PathPoint[])[(path as PathPoint[]).length - 1] : null;

    return (
        <MapContainer
            center={center}
            zoom={14}
            scrollWheelZoom={false}
            className="h-full w-full"
            style={{ minHeight: "100%" }}
        >
            <TileLayer
                attribution='&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
                url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            />

            {hasPath ? (
                <>
                    <FitBounds points={pathCoords} />

                    <Polyline
                        positions={pathCoords}
                        pathOptions={{
                            color: "#3b82f6",
                            weight: 3,
                            opacity: 0.75,
                            dashArray: "8 6",
                        }}
                    />

                    {userPoint && (
                        <Marker
                            position={[userPoint.lat, userPoint.lon]}
                            icon={userIcon}
                        >
                            <Popup>
                                <span className="text-xs font-medium">Your location</span>
                            </Popup>
                        </Marker>
                    )}

                    {eventPoint && (
                        <Marker
                            position={[eventPoint.lat, eventPoint.lon]}
                            icon={eventIcon}
                        >
                            <Popup>
                                <span className="text-xs font-medium">{eventPoint.name}</span>
                            </Popup>
                        </Marker>
                    )}
                </>
            ) : (
                <>
                    <ChangeCenter center={center} />
                    <Marker position={center} icon={eventIcon}>
                        <Popup>
                            <span className="text-xs font-medium">Event location</span>
                        </Popup>
                    </Marker>
                </>
            )}
        </MapContainer>
    );
}