import { MapContainer, Marker, Popup, TileLayer } from "react-leaflet";
import { ChangeMapView } from "../features/events/utils/map";
import type { locationResp } from "../features/events/types/locationResponse";
import L from "leaflet";
import markerIcon2x from "leaflet/dist/images/marker-icon-2x.png";
import markerIcon from "leaflet/dist/images/marker-icon.png";
import markerShadow from "leaflet/dist/images/marker-shadow.png";

interface Props {
    position: [number, number];
    selectedLocation: locationResp | undefined;
};

delete (L.Icon.Default.prototype as any)._getIconUrl;
L.Icon.Default.mergeOptions({
    iconUrl: markerIcon,
    iconRetinaUrl: markerIcon2x,
    shadowUrl: markerShadow,
});

export default function MapPicker({ position, selectedLocation }: Props) {
    return (
        <div className="relative h-60 sm:h-72 lg:h-80 w-full rounded-xl overflow-hidden">
            <MapContainer
                center={position}
                zoom={13}
                scrollWheelZoom={false}
                className="h-full w-full"
            >
                <TileLayer
                    attribution='&copy; OpenStreetMap contributors'
                    url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
                />

                <ChangeMapView position={position} />

                <Marker position={position}>
                    <Popup>
                        {selectedLocation?.display_name || "Selected location"}
                    </Popup>
                </Marker>
            </MapContainer>
        </div>
    )
}