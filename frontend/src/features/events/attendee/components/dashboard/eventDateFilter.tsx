import { createPortal } from "react-dom";
import { FiCalendar, FiChevronDown } from "react-icons/fi";
import { IoClose } from "react-icons/io5";
import type { DateFilterType } from "../../../../../utils/dateFilter";
import DateFilterDropdown from "../../../../../components/dateFilter/dateFilterDropdown";
import { useBottomSheetOrDropdown } from "../../../../../utils/useBottomSheetorDropdown";

interface Props {
    dateFilterType: DateFilterType;
    setDateFilterType: React.Dispatch<React.SetStateAction<DateFilterType>>;
    selectedDate: Date | null;
    setSelectedDate: React.Dispatch<React.SetStateAction<Date | null>>;
    selectedMonth?: number;
    setSelectedMonth: React.Dispatch<React.SetStateAction<number | undefined>>;
    selectedYear?: number;
    setSelectedYear: React.Dispatch<React.SetStateAction<number | undefined>>;
}

const DATE_FILTER_LABELS: Record<DateFilterType, string> = {
    all: "Any Time",
    today: "Today",
    this_month: "This Month",
    this_year: "This Year",
    date: "Select Date",
    month: "Select Month",
    year: "Select Year",
};

const NO_EXTRA_INPUT: DateFilterType[] = ["all", "today", "this_month", "this_year"];

export default function EventDateFilter({
    dateFilterType,
    setDateFilterType,
    selectedDate,
    setSelectedDate,
    selectedMonth,
    setSelectedMonth,
    selectedYear,
    setSelectedYear,
}: Props) {
    const { open, isMobile, handleToggle, close, buttonRef, dropdownRef, dropdownPos } = useBottomSheetOrDropdown();

    const isActive = dateFilterType !== "all";
    const label = DATE_FILTER_LABELS[dateFilterType] ?? "Any Time";

    const handleSetDateFilterType = (val: React.SetStateAction<DateFilterType>) => {
        setDateFilterType(val);
        if (isMobile && typeof val === "string" && NO_EXTRA_INPUT.includes(val)) {
            setTimeout(close, 150);
        }
    };

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
                <FiCalendar className={`w-4 h-4 ${isActive ? "text-blue-500" : "text-slate-400"}`} />
                <span className="flex-1 text-left">{label}</span>
                <FiChevronDown
                    className={`w-4 h-4 text-slate-400 transition-transform ${open ? "rotate-180" : ""}`}
                />
            </button>

            {!isMobile && open && createPortal(
                <div
                    ref={dropdownRef}
                    style={{ position: "fixed", top: dropdownPos.top, left: dropdownPos.left, zIndex: 9999 }}
                    className="bg-white border border-slate-200 rounded-xl shadow-xl w-[320px] p-4"
                >
                    <DateFilterDropdown
                        dateFilterType={dateFilterType}
                        setDateFilterType={setDateFilterType}
                        selectedDate={selectedDate}
                        setSelectedDate={setSelectedDate}
                        selectedMonth={selectedMonth}
                        setSelectedMonth={setSelectedMonth}
                        selectedYear={selectedYear}
                        setSelectedYear={setSelectedYear}
                    />
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
                            <h2 className="text-base font-semibold text-gray-800">Filter by Date</h2>
                            <button
                                onClick={close}
                                className="p-1.5 rounded-lg hover:bg-gray-100 text-gray-500 transition-colors"
                            >
                                <IoClose size={20} />
                            </button>
                        </div>

                        <div className="px-5 py-4 pb-8">
                            <DateFilterDropdown
                                dateFilterType={dateFilterType}
                                setDateFilterType={handleSetDateFilterType}
                                selectedDate={selectedDate}
                                setSelectedDate={setSelectedDate}
                                selectedMonth={selectedMonth}
                                setSelectedMonth={setSelectedMonth}
                                selectedYear={selectedYear}
                                setSelectedYear={setSelectedYear}
                            />
                        </div>
                    </div>
                </>,
                document.body
            )}
        </>
    );
}