import { FiNavigation } from "react-icons/fi";
import type { PathPoint } from "../../../types/publicEventResponse";


interface Props {
    address: string;
    detailAddress: string;
    coordinates: { lat: number; lon: number };
    path?: PathPoint[];
    distance?: string;
    MapComponent: React.ComponentType<MapComponentProps>;
}

export interface MapComponentProps {
    center: [number, number];
    path?: PathPoint[];
}

export default function LocationCard({
    address,
    detailAddress,
    coordinates,
    path,
    MapComponent,
}: Props) {
    const handleSeeRoute = () => {
        const dest = `${coordinates.lat},${coordinates.lon}`;
        window.open(`https://www.google.com/maps/dir/?api=1&destination=${dest}`, "_blank");
    };

    return (
        <div className="bg-white rounded-2xl border border-slate-200 p-6 space-y-4">
            <h2 className="text-lg font-bold text-slate-900">Location</h2>

            {/* Map */}
            <div className="rounded-xl overflow-hidden border border-slate-100 h-56">
                <MapComponent
                    center={[coordinates.lat, coordinates.lon]}
                    path={path}
                />
            </div>

            {/* Address text */}
            <div>
                <p className="text-sm font-semibold text-slate-800">{detailAddress}</p>
                <p className="text-xs text-slate-500 mt-1 leading-relaxed">{address}</p>
            </div>

            {/* See Route button */}
            <button
                onClick={handleSeeRoute}
                className="w-full flex items-center justify-center gap-2 bg-slate-900 hover:bg-slate-700 text-white text-sm font-medium py-3 rounded-xl transition-colors"
            >
                <FiNavigation className="w-4 h-4" />
                See Route
            </button>
        </div>
    );
}