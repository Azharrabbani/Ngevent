import { useEffect, useRef } from "react";
import L from "leaflet";
import "leaflet/dist/leaflet.css";
import type { MapComponentProps } from "../cards/locationCard";

delete (L.Icon.Default.prototype as any)._getIconUrl;
L.Icon.Default.mergeOptions({
    iconUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png",
    iconRetinaUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png",
    shadowUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png",
});

const EVENT_ICON = L.icon({
    iconUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png",
    iconRetinaUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png",
    shadowUrl: "https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png",
    iconSize: [25, 41],
    iconAnchor: [12, 41],
    popupAnchor: [1, -34],
    shadowSize: [41, 41],
});

const USER_ICON = L.divIcon({
    className: "",
    html: `<div style="
        width: 14px; height: 14px;
        background: #3b82f6;
        border: 3px solid white;
        border-radius: 50%;
        box-shadow: 0 0 0 2px #3b82f6;
    "></div>`,
    iconSize: [14, 14],
    iconAnchor: [7, 7],
});

export default function EventMap({ center, path }: MapComponentProps) {
    const containerRef = useRef<HTMLDivElement>(null);
    const mapRef = useRef<L.Map | null>(null);
    const polylineRef = useRef<L.Polyline | null>(null);
    const userMarkerRef = useRef<L.Marker | null>(null);
    const eventMarkerRef = useRef<L.Marker | null>(null);

    // Initialize map
    useEffect(() => {
        if (!containerRef.current || mapRef.current) return;

        const map = L.map(containerRef.current, {
            center,
            zoom: 15,
            zoomControl: true,
            scrollWheelZoom: false,
        });

        L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
            attribution: '© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>',
            maxZoom: 19,
        }).addTo(map);

        // Event location marker (always present)
        eventMarkerRef.current = L.marker(center, { icon: EVENT_ICON })
            .addTo(map);

        mapRef.current = map;

        return () => {
            map.remove();
            mapRef.current = null;
        };
    }, []);

    // Draw/clear polyline and user marker
    useEffect(() => {
        const map = mapRef.current;
        if (!map) return;

        // Clear previous route layers
        polylineRef.current?.remove();
        userMarkerRef.current?.remove();
        polylineRef.current = null;
        userMarkerRef.current = null;

        if (!path || path.length === 0) {
            // show the event marker at current center when there is no path
            map.setView(center, 15, { animate: true });
            return;
        }

        const latLngs: L.LatLngExpression[] = path.map((p) => [p.lat, p.lon]);

        // Draw the route polyline
        polylineRef.current = L.polyline(latLngs, {
            color: "#3b82f6",
            weight: 4,
            opacity: 0.85,
            lineJoin: "round",
            lineCap: "round",
        }).addTo(map);

        // User position marker at the first point
        const userPoint = path[0];
        userMarkerRef.current = L.marker([userPoint.lat, userPoint.lon], {
            icon: USER_ICON,
        })
            .bindTooltip("You", { permanent: false, direction: "top" })
            .addTo(map);

        // Fit the map to show the full route with padding
        map.fitBounds(polylineRef.current.getBounds(), {
            padding: [40, 40],
            animate: true,
            duration: 0.6,
        });
    }, [path, center]);

    // Keep event marker in sync if center ever changes
    useEffect(() => {
        eventMarkerRef.current?.setLatLng(center);
    }, [center]);

    return (
        <div
            ref={containerRef}
            style={{ width: "100%", height: "100%" }}
        />
    );
}