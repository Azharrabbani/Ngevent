import { IoFilter, IoClose } from "react-icons/io5";
import Input from "../../../../../components/input";
import { SearchIcon } from "../../../../../components/icon";
import type { DateFilterType } from "../../../../../utils/dateFilter";
import type { Dispatch, SetStateAction } from "react";
import EventCreateDateFilterDropdown from "../dateFilterDropdown";
import DateFilterDropdown from "../../../../../components/dateFilter/dateFilterDropdown";
import { createPortal } from "react-dom";
import { useBottomSheetOrDropdown } from "../../../../../utils/useBottomSheetorDropdown";

interface Props {
    reviewEvent: boolean;
    status: string;
    search?: string;
    setSearch?: (val: string | undefined) => void;
    sort?: string;
    setSort?: React.Dispatch<React.SetStateAction<string | undefined>>;
    dateFilter?: string;
    setDateFilter?: React.Dispatch<React.SetStateAction<string | undefined>>;
    dateFilterType: DateFilterType;
    setDateFilterType: Dispatch<SetStateAction<DateFilterType>>;
    selectedDate: Date | null;
    setSelectedDate: Dispatch<SetStateAction<Date | null>>;
    selectedMonth?: number;
    setSelectedMonth: Dispatch<SetStateAction<number | undefined>>;
    selectedYear?: number;
    setSelectedYear: Dispatch<SetStateAction<number | undefined>>;
    getUpdate?: boolean;
    setGetupdate?: (val: boolean | undefined) => void;
}

const NO_EXTRA_INPUT = ["all", "today", "this_month", "this_year"];

export default function EventsHeader({
    reviewEvent,
    status,
    search,
    setSearch,
    sort,
    setSort,
    dateFilter,
    setDateFilter,
    dateFilterType,
    setDateFilterType,
    selectedDate,
    setSelectedDate,
    selectedMonth,
    setSelectedMonth,
    selectedYear,
    setSelectedYear,
    getUpdate,
    setGetupdate,
}: Props) {
    const { open, isMobile, handleToggle, close, buttonRef, dropdownRef, dropdownPos } =
        useBottomSheetOrDropdown();

    const activeFilterLabel = (() => {
        if (dateFilterType === "all") return null;
        const labels: Record<DateFilterType, string> = {
            all: "",
            today: "Today",
            this_month: "This Month",
            this_year: "This Year",
            date: "Specific Date",
            month: "Month",
            year: "Year",
        };
        return labels[dateFilterType];
    })();

    const handleSetDateFilterType = (val: React.SetStateAction<DateFilterType>) => {
        setDateFilterType(val);
        if (isMobile && typeof val === "string" && NO_EXTRA_INPUT.includes(val)) {
            setTimeout(close, 150);
        }
    };

    return (
        <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4 p-6 border-b border-gray-100">
            <div>
                <h1 className="text-xl font-semibold tracking-wide text-gray-700 uppercase">
                    {reviewEvent ? "Event Submissions" : `${status} Events`}
                </h1>
            </div>

            <div className="flex items-center gap-3 flex-wrap">
                <div className="relative w-full lg:w-[280px]">
                    <Input
                        leftIcon={<SearchIcon />}
                        type="text"
                        placeholder="Search events, organizers..."
                        value={search}
                        onChange={(e) => setSearch?.(e.target.value)}
                        className="w-full bg-white pl-10 pr-4 py-3 rounded-lg border border-gray-300 outline-none focus:ring-2 focus:ring-blue-500 text-sm"
                    />
                </div>

                {reviewEvent && (
                    <button
                        onClick={() => setGetupdate?.(getUpdate === true ? false : true)}
                        className={`px-4 py-2 rounded-lg border text-sm font-medium transition-all ${getUpdate === true
                            ? "bg-blue-50 border-blue-500 text-blue-600"
                            : "border-gray-300 text-gray-500 hover:bg-gray-50"
                            }`}
                    >
                        Update requests
                    </button>
                )}

                {sort && (
                    <button
                        onClick={() => setSort?.(sort === "desc" ? "asc" : "desc")}
                        className="flex items-center gap-2 px-4 py-2 rounded-lg border border-gray-300 text-sm text-gray-500 hover:bg-gray-50"
                    >
                        <IoFilter size={16} />
                        {sort === "desc" ? "Oldest first" : "Newest first"}
                    </button>
                )}

                <EventCreateDateFilterDropdown
                    dateFilter={dateFilter}
                    setDateFilter={setDateFilter}
                />

                <button
                    ref={buttonRef}
                    onClick={handleToggle}
                    className={`flex items-center gap-2 px-4 py-2 rounded-lg border text-sm transition-all ${open || activeFilterLabel
                        ? "bg-blue-50 border-blue-500 text-blue-600"
                        : "border-gray-300 text-gray-500 hover:bg-gray-50"
                        }`}
                >
                    {activeFilterLabel ?? "Event date"}
                    <IoFilter size={14} />
                </button>
            </div>

            {!isMobile && open && createPortal(
                <div
                    ref={dropdownRef}
                    style={{
                        position: "fixed",
                        top: dropdownPos.top,
                        left: dropdownPos.left,
                        zIndex: 9999,
                    }}
                    className="bg-white border border-gray-200 rounded-xl shadow-lg p-3 w-[220px]"
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
                            <h2 className="text-base font-semibold text-gray-800">
                                Filter by Event Date
                            </h2>
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
        </div>
    );
}