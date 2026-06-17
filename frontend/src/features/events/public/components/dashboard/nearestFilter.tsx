import { PinIcon, SpinnerIcon } from "../../../../../components/icon";

interface Props {
    enabled: boolean;
    onToggle: (value: boolean) => void;
    locationLoading: boolean;
    denied: boolean;
}

export default function NearestFilter({ enabled, onToggle, locationLoading, denied }: Props) {
    const handleClick = () => {
        if (denied || locationLoading) return;
        onToggle(!enabled);
    };

    const isDenied = denied && !locationLoading;

    return (
        <div className="relative group">
            <button
                type="button"
                onClick={handleClick}
                disabled={isDenied || locationLoading}
                className={[
                    "inline-flex items-center gap-2 px-4 py-2 rounded-full border text-sm font-medium transition-all duration-200",
                    "focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2",
                    isDenied
                        ? "border-slate-200 bg-slate-50 text-slate-400 cursor-not-allowed"
                        : locationLoading
                            ? "border-blue-200 bg-blue-50 text-blue-400 cursor-wait"
                            : enabled
                                ? "border-blue-500 bg-blue-500 text-white shadow-sm shadow-blue-200"
                                : "border-slate-200 bg-white text-slate-600 hover:border-blue-300 hover:text-blue-600 hover:bg-blue-50",
                ].join(" ")}
                aria-pressed={enabled}
                aria-label={
                    isDenied
                        ? "Location access denied"
                        : enabled
                            ? "Disable nearest filter"
                            : "Show nearest events"
                }
            >
                {locationLoading ? (
                    <SpinnerIcon className="w-4 h-4 animate-spin" />
                ) : isDenied ? (
                    <PinIcon className="w-4 h-4" />
                ) : (
                    <PinIcon className={`w-4 h-4 ${enabled ? "fill-white/20" : ""}`} />
                )}

                <span>
                    {locationLoading
                        ? "Locating..."
                        : isDenied
                            ? "Location denied"
                            : "Nearest"}
                </span>
            </button>

            {isDenied && (
                <div
                    role="tooltip"
                    className={[
                        "absolute bottom-full left-1/2 -translate-x-1/2 mb-2 w-56 z-10",
                        "px-3 py-2 rounded-lg bg-slate-800 text-white text-xs text-center leading-snug",
                        "opacity-0 group-hover:opacity-100 pointer-events-none transition-opacity duration-150",
                        "after:content-[''] after:absolute after:top-full after:left-1/2 after:-translate-x-1/2",
                        "after:border-4 after:border-transparent after:border-t-slate-800",
                    ].join(" ")}
                >
                    Enable location access in your browser settings to use this filter.
                </div>
            )}
        </div>
    );
}