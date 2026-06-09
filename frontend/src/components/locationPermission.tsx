import { FiMapPin } from "react-icons/fi";

interface Props {
    onRequestLocation: () => void;
}

export default function LocationPermissionBanner({ onRequestLocation }: Props) {
    return (
        <div className="flex items-center gap-3 bg-blue-50 border border-blue-200 rounded-xl px-4 py-3 text-sm">
            <FiMapPin className="w-4 h-4 text-blue-500 shrink-0" />
            <p className="text-blue-700 flex-1">
                Enable location to see how far this event is from you and to visualize the fastest route.
            </p>
            <button
                onClick={onRequestLocation}
                className="text-blue-600 font-semibold hover:underline text-xs shrink-0"
            >
                Enable
            </button>
        </div>
    );
}