import { MdOutlineCalendarMonth, MdOutlineCalendarToday } from "react-icons/md";

interface Props {
    value:    "monthly" | "yearly";
    onChange: (value: "monthly" | "yearly") => void;
}

export default function PeriodToggle({ value, onChange }: Props) {
    return (
        <div className="mb-6">
            <label className="block text-sm font-medium text-gray-700 mb-2">
                Period Type
            </label>
            <div className="inline-flex rounded-xl bg-gray-100 p-1 gap-1">
                <button
                    onClick={() => onChange("monthly")}
                    className={`flex items-center gap-2 px-5 py-2 rounded-lg text-sm font-semibold transition-all duration-200
                        ${value === "monthly"
                            ? "bg-[#0056D2] text-white shadow"
                            : "text-gray-500 hover:text-gray-700"}`}
                >
                    <MdOutlineCalendarMonth className="text-base" />
                    Monthly
                </button>
                <button
                    onClick={() => onChange("yearly")}
                    className={`flex items-center gap-2 px-5 py-2 rounded-lg text-sm font-semibold transition-all duration-200
                        ${value === "yearly"
                            ? "bg-[#0056D2] text-white shadow"
                            : "text-gray-500 hover:text-gray-700"}`}
                >
                    <MdOutlineCalendarToday className="text-base" />
                    Yearly
                </button>
            </div>
        </div>
    );
}