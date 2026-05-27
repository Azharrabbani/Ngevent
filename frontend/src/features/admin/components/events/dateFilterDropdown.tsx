import { useRef, useState, useEffect } from "react";
import { FiCalendar } from "react-icons/fi";
import { IoCheckmark } from "react-icons/io5";

const options = [
    { label: "Last 7 days", value: "week" },
    { label: "Last 30 days", value: "month" },
    { label: "Last year", value: "year" },
];

interface Props {
    dateFilter?: string;
    setDateFilter?: (val: string | undefined) => void;
}

export default function DateFilterDropdown({ dateFilter, setDateFilter }: Props) {
    const [open, setOpen] = useState(false);
    const ref = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const handler = (e: MouseEvent) => {
            if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
        };
        document.addEventListener("mousedown", handler);
        return () => document.removeEventListener("mousedown", handler);
    }, []);

    const selected = options.find((o) => o.value === dateFilter);

    return (
        <div ref={ref} className="relative">
            <button
                onClick={() => setOpen(!open)}
                className={`flex items-center gap-2 px-4 py-2 rounded-lg border text-sm transition-all ${dateFilter
                    ? "border-blue-500 text-blue-600 bg-blue-50"
                    : "border-gray-300 text-gray-500 hover:bg-gray-50"
                    }`}
            >
                <FiCalendar size={15} />
                {selected?.label ?? "Date range"}
            </button>

            {open && (
                <div className="absolute right-0 mt-2 w-44 bg-white border border-gray-200 rounded-xl shadow-md z-50 p-1.5">
                    {options.map((opt) => (
                        <button
                            key={opt.value}
                            onClick={() => {
                                setDateFilter?.(dateFilter === opt.value ? undefined : opt.value);
                                setOpen(false);
                            }}
                            className="flex justify-between items-center w-full text-left px-3 py-2 rounded-lg text-sm hover:bg-blue-500 hover:text-white transition group"
                        >
                            {opt.label}
                            {dateFilter === opt.value && <IoCheckmark size={15} />}
                        </button>
                    ))}
                </div>
            )}
        </div>
    );
}