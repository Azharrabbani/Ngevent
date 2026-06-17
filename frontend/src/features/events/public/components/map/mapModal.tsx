import { useEffect } from "react";
import { createPortal } from "react-dom";
import type { MapComponentProps } from "../cards/locationCard";
import type { PathPoint } from "../../../types/publicEventResponse";
import { CloseIcon } from "../../../../../components/icon";

interface Props {
    isOpen: boolean;
    onClose: () => void;
    center: [number, number];
    path?: PathPoint[];
    MapComponent: React.ComponentType<MapComponentProps>;
    title?: string;
}

export default function MapModal({
    isOpen,
    onClose,
    center,
    path,
    MapComponent,
    title,
}: Props) {
    useEffect(() => {
        if (!isOpen) return;
        const handler = (e: KeyboardEvent) => {
            if (e.key === "Escape") onClose();
        };
        document.addEventListener("keydown", handler);
        return () => document.removeEventListener("keydown", handler);
    }, [isOpen, onClose]);

    useEffect(() => {
        if (isOpen) {
            document.body.style.overflow = "hidden";
        } else {
            document.body.style.overflow = "";
        }
        return () => {
            document.body.style.overflow = "";
        };
    }, [isOpen]);

    if (!isOpen) return null;

    return createPortal(
        <div
            className="fixed inset-0 z-50 flex flex-col"
            role="dialog"
            aria-modal="true"
            aria-label={title ? `Map — ${title}` : "Map"}
        >
            <div
                className="absolute inset-0 bg-black/60 backdrop-blur-sm"
                onClick={onClose}
            />
            <div className="relative flex flex-col w-full h-full max-w-5xl mx-auto my-6 px-4 sm:px-6">
                <div className="relative flex flex-col flex-1 bg-white rounded-2xl overflow-hidden shadow-2xl">

                    <div className="flex items-center justify-between px-5 py-3 border-b border-slate-100 shrink-0">
                        <p className="text-sm font-semibold text-slate-800 truncate">
                            {title ?? "Location"}
                        </p>
                        <button
                            onClick={onClose}
                            className="p-1.5 rounded-lg text-slate-500 hover:text-slate-800 hover:bg-slate-100 transition-colors focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-500"
                            aria-label="Close map"
                        >
                            <CloseIcon className="w-4 h-4" />
                        </button>
                    </div>

                    <div className="flex-1 min-h-0">
                        <MapComponent center={center} path={path} />
                    </div>
                </div>
            </div>
        </div>,
        document.body
    );
}