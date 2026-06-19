import { HiOutlineRefresh } from "react-icons/hi";
import PeriodToggle from "./reportToggle";

const MONTHS = [
    "January", "February", "March", "April",
    "May", "June", "July", "August",
    "September", "October", "November", "December",
];

const currentYear = new Date().getFullYear();
const YEAR_OPTIONS = Array.from({ length: 6 }, (_, i) => currentYear - i);

interface Props {
    period: "monthly" | "yearly";
    month: number;
    year: number;
    isLoading: boolean;
    onPeriodChange: (p: "monthly" | "yearly") => void;
    onMonthChange: (m: number) => void;
    onYearChange: (y: number) => void;
    onGenerate: () => void;
}

export default function ReportFilters({
    period,
    month,
    year,
    isLoading,
    onPeriodChange,
    onMonthChange,
    onYearChange,
    onGenerate,
}: Props) {
    return (
        <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-6 mb-6">
            <p className="text-xs font-semibold text-gray-400 uppercase tracking-widest mb-5">
                Report Settings
            </p>

            <PeriodToggle value={period} onChange={onPeriodChange} />

            <div className="flex flex-wrap gap-4 items-end">
                {period === "monthly" && (
                    <div className="flex flex-col gap-1">
                        <label className="text-sm font-medium text-gray-700">Month</label>
                        <select
                            value={month}
                            onChange={(e) => onMonthChange(Number(e.target.value))}
                            className="border border-gray-200 rounded-xl px-4 py-2.5 text-sm text-gray-800 bg-white focus:outline-none focus:ring-2 focus:ring-[#0056D2]/30 focus:border-[#0056D2] transition min-w-[160px]"
                        >
                            {MONTHS.map((m, i) => (
                                <option key={i} value={i + 1}>{m}</option>
                            ))}
                        </select>
                    </div>
                )}

                <div className="flex flex-col gap-1">
                    <label className="text-sm font-medium text-gray-700">Year</label>
                    <select
                        value={year}
                        onChange={(e) => onYearChange(Number(e.target.value))}
                        className="border border-gray-200 rounded-xl px-4 py-2.5 text-sm text-gray-800 bg-white focus:outline-none focus:ring-2 focus:ring-[#0056D2]/30 focus:border-[#0056D2] transition min-w-[120px]"
                    >
                        {YEAR_OPTIONS.map((y) => (
                            <option key={y} value={y}>{y}</option>
                        ))}
                    </select>
                </div>

                <button
                    onClick={onGenerate}
                    disabled={isLoading}
                    className="flex items-center gap-2 bg-[#0056D2] hover:bg-[#0046b0] disabled:bg-[#0056D2]/50 text-white text-sm font-semibold px-6 py-2.5 rounded-xl transition-all duration-200 shadow-sm"
                >
                    {isLoading ? (
                        <>
                            <svg className="animate-spin h-4 w-4" viewBox="0 0 24 24" fill="none">
                                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z" />
                            </svg>
                            Generating...
                        </>
                    ) : (
                        <>
                            <HiOutlineRefresh className="text-base" />
                            Generate Report
                        </>
                    )}
                </button>
            </div>
        </div>
    );
}