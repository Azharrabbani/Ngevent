import { createPortal } from "react-dom";
import { FiMapPin, FiChevronDown, FiCheck } from "react-icons/fi";
import { IoClose } from "react-icons/io5";
import { LOCATIONS } from "../../constants/location";
import { useBottomSheetOrDropdown } from "../../../../../utils/useBottomSheetorDropdown";

interface Props {
    value?: string;
    onChange: (value?: string) => void;
}

export default function LocationFilter({ value, onChange }: Props) {
    const { open, isMobile, handleToggle, close, buttonRef, dropdownRef, dropdownPos } = useBottomSheetOrDropdown();

    const handleSelect = (location: string) => {
        onChange(location === "All Locations" ? undefined : location);
        close();
    };

    const isActive = !!value;

    return (
        <>
            <button
                ref={buttonRef}
                onClick={handleToggle}
                className={`
                    flex items-center gap-2 border rounded-xl px-4 py-2.5
                    bg-white min-w-[200px] text-sm transition
                    ${isActive || open
                        ? "border-blue-400 text-blue-700 font-medium"
                        : "border-slate-200 text-slate-600 hover:border-slate-300"
                    }
                `}
            >
                <FiMapPin className={`w-4 h-4 ${isActive ? "text-blue-500" : "text-slate-400"}`} />
                <span className="flex-1 text-left">{value ?? "All Locations"}</span>
                <FiChevronDown
                    className={`w-4 h-4 text-slate-400 transition-transform ${open ? "rotate-180" : ""}`}
                />
            </button>

            {!isMobile && open && createPortal(
                <div
                    ref={dropdownRef}
                    style={{ position: "fixed", top: dropdownPos.top, left: dropdownPos.left, zIndex: 9999 }}
                    className="bg-white border border-slate-200 rounded-xl shadow-lg overflow-hidden py-1 w-[220px]"
                >
                    {LOCATIONS.map((location) => (
                        <button
                            key={location}
                            onClick={() => handleSelect(location)}
                            className={`
                                w-full text-left px-4 py-2.5 text-sm flex items-center justify-between transition
                                ${location === (value ?? "All Locations")
                                    ? "text-blue-600 font-medium bg-blue-50"
                                    : "text-slate-700 hover:bg-slate-50"
                                }
                            `}
                        >
                            {location}
                            {location === (value ?? "All Locations") && (
                                <FiCheck className="w-4 h-4 text-blue-500" />
                            )}
                        </button>
                    ))}
                </div>,
                document.body
            )}

            {isMobile && createPortal(
                <>
                    <div
                        onClick={close}
                        className={`fixed inset-0 bg-black/40 z-[9998] transition-opacity duration-300 ${open ? "opacity-100" : "opacity-0 pointer-events-none"
                            }`}
                    />

                    <div
                        className={`fixed bottom-0 left-0 right-0 z-[9999] bg-white rounded-t-2xl shadow-2xl
                            transition-transform duration-300 ease-out
                            ${open ? "translate-y-0" : "translate-y-full"}`}
                    >
                        <div className="flex justify-center pt-3 pb-1">
                            <div className="w-10 h-1 rounded-full bg-gray-300" />
                        </div>
                        <div className="flex items-center justify-between px-5 py-3 border-b border-gray-100">
                            <h2 className="text-base font-semibold text-gray-800">Select Location</h2>
                            <button
                                onClick={close}
                                className="p-1.5 rounded-lg hover:bg-gray-100 text-gray-500 transition-colors"
                            >
                                <IoClose size={20} />
                            </button>
                        </div>

                        <div className="overflow-y-auto max-h-[60vh] py-2 pb-8">
                            {LOCATIONS.map((location) => (
                                <button
                                    key={location}
                                    onClick={() => handleSelect(location)}
                                    className={`
                                        w-full text-left px-5 py-3.5 text-sm flex items-center justify-between transition
                                        ${location === (value ?? "All Locations")
                                            ? "text-blue-600 font-medium bg-blue-50"
                                            : "text-slate-700 active:bg-slate-50"
                                        }
                                    `}
                                >
                                    {location}
                                    {location === (value ?? "All Locations") && (
                                        <FiCheck className="w-4 h-4 text-blue-500" />
                                    )}
                                </button>
                            ))}
                        </div>
                    </div>
                </>,
                document.body
            )}
        </>
    );
}