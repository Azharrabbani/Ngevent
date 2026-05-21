import Pagination from "../../../../components/pagination";
import EventTable from "./eventTable";
import EventCard from "./eventCard";
import type { EventsResponse } from "../../../events/types/organizerResponse";
import type { PaginatedData } from "../../../../types/apiResponse";
import { FiSearch } from "react-icons/fi";
import { IoFilter } from "react-icons/io5";
import DateFilterDropdown from "./dateFilterDropdown";


interface Props {
    status: string;
    data: PaginatedData<EventsResponse> | undefined;
    isLoading: boolean;
    reviewEvent: boolean;
    currentPage: number;
    totalPage: number;
    search?: string;
    setSearch?: (val: string | undefined) => void;
    setCurrentPage: React.Dispatch<React.SetStateAction<number>>;
    sort?: string;
    setSort?: React.Dispatch<React.SetStateAction<string | undefined>>;
    dateFilter?: string;
    setDateFilter?: React.Dispatch<React.SetStateAction<string | undefined>>;
    getUpdate?: boolean;
    setGetupdate?: React.Dispatch<React.SetStateAction<boolean | undefined>>;
}
export default function EventList({
    status,
    data,
    isLoading,
    reviewEvent,
    currentPage,
    totalPage,
    search,
    setSearch,
    setCurrentPage,
    sort,
    setSort,
    dateFilter,
    setDateFilter,
    getUpdate,
    setGetupdate }: Props) {
    return (
        <div className="bg-white rounded-3xl shadow-sm border border-gray-100 overflow-hidden">

            {/* Header */}
            <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4 p-6 border-b border-gray-100">

                <div>
                    <h1 className="text-xl font-semibold tracking-wide text-gray-700 uppercase">
                        {reviewEvent ? "Event Submissions" : `${status} Events`}
                    </h1>
                </div>

                <div className="flex items-center gap-3 flex-wrap">

                    {/* Search */}
                    <div className="relative w-full lg:w-[280px]">
                        <FiSearch
                            className="absolute left-4 top-1/2 -translate-y-1/2 text-gray-400"
                            size={18}
                        />

                        <input
                            type="text"
                            placeholder="Search events, organizers..."
                            value={search}
                            onChange={(e) => setSearch?.(e.target.value)}
                            className="w-full pl-10 pr-4 py-2.5 rounded-lg border border-gray-300 outline-none focus:ring-2 focus:ring-blue-500 text-sm"
                        />
                    </div>

                    {/* Update Request Filter */}
                    {reviewEvent && (
                        <button
                            onClick={() =>
                                setGetupdate?.(
                                    getUpdate === true ? undefined : true
                                )
                            }
                            className={`px-4 py-2 rounded-lg border text-sm font-medium transition-all ${getUpdate === true
                                ? "bg-blue-50 border-blue-500 text-blue-600"
                                : "border-gray-300 text-gray-500 hover:bg-gray-50"
                                }`}
                        >
                            Update requests
                        </button>
                    )}

                    {/* Sort */}
                    <button
                        onClick={() =>
                            setSort?.(sort === "desc" ? "asc" : "desc")
                        }
                        className="flex items-center gap-2 px-4 py-2 rounded-lg border border-gray-300 text-sm text-gray-500 hover:bg-gray-50"
                    >
                        <IoFilter size={16} />

                        {sort === "desc"
                            ? "Newest first"
                            : "Oldest first"}
                    </button>

                    {/* Date Filter */}
                    <DateFilterDropdown
                        dateFilter={dateFilter}
                        setDateFilter={setDateFilter}
                    />
                </div>
            </div>

            {/* Desktop */}
            <EventTable
                status={status}
                data={data}
                isLoading={isLoading}
                isReview={reviewEvent}
                sort={sort}
                setSort={setSort}
            />

            {/* Mobile */}
            <EventCard
                status={status}
                data={data}
                isLoading={isLoading}
                isReview={reviewEvent}
                sort={sort}
                setSort={setSort}
            />

            {data && (
                <div className="flex flex-col md:flex-row items-center justify-between gap-4 px-6 md:px-8 border-t border-gray-100">

                    <p className="text-sm text-gray-500">
                        Showing 1 to 4 of {data?.total_rows} events
                    </p>

                    <Pagination
                        currentPage={currentPage}
                        totalPage={totalPage}
                        onPrev={() =>
                            setCurrentPage((prev) =>
                                Math.max(prev - 1, 1)
                            )
                        }
                        onNext={() =>
                            setCurrentPage((prev) =>
                                Math.min(prev + 1, totalPage)
                            )
                        }
                        onCurrent={(page) => setCurrentPage(page)}
                    />
                </div>
            )}
        </div>
    )
}