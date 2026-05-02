import { MapContainer, Marker, Popup, TileLayer } from "react-leaflet";
import { ChangeMapView } from "../features/events/utils/map";
import type { locationResp } from "../features/events/types/locationResponse";

interface Props {
    position: [number, number];
    selectedLocation: locationResp | undefined;
};

export default function MapPicker({ position, selectedLocation }: Props) {
    return (
        <div className="h-60 z-10 sm:h-72 lg:h-80 w-full rounded-xl overflow-hidden">
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