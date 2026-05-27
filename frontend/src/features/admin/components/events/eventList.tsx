import Pagination from "../../../../components/pagination";
import EventTable from "./eventTable";
import EventCard from "./eventCard";
import type { EventsResponse } from "../../../events/types/eventResponse";
import type { PaginatedData } from "../../../../types/apiResponse";
import EventsHeader from "./header";
import { SpinnerIcon } from "../../../../components/icon";


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
    setGetupdate?: (val: boolean | undefined) => void;
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
    const isEmpty = !data || data.total_rows === 0;

    return (
        <div className="bg-white rounded-3xl shadow-sm border border-gray-100 overflow-hidden">
            <EventsHeader
                reviewEvent={reviewEvent}
                status={status}
                search={search}
                setSearch={setSearch}
                sort={sort}
                setSort={setSort}
                dateFilter={dateFilter}
                setDateFilter={setDateFilter}
                getUpdate={getUpdate}
                setGetupdate={setGetupdate}
            />

            {isLoading ? (
                <div className="flex justify-center py-20">
                    <SpinnerIcon className="animate-spin w-8 h-8 text-blue-500" />
                </div>
            ) : isEmpty ? (
                <div className="flex flex-col items-center justify-center py-20 text-center">
                    <h1 className="text-xl font-semibold text-gray-700">
                        No events found
                    </h1>
                </div>
            ) : (
                <>
                    <EventTable
                        key={`table-${getUpdate ?? "false"}`}
                        status={status}
                        data={data}
                        isReview={reviewEvent}
                        sort={sort}
                        setSort={setSort}
                    />

                    <EventCard
                        key={`card-${getUpdate ?? "false"}`}
                        status={status}
                        data={data}
                        isReview={reviewEvent}
                        sort={sort}
                        setSort={setSort}
                    />

                    <div className="flex flex-col md:flex-row items-center justify-between gap-4 px-6 md:px-8 border-t border-gray-100">

                        <p className="text-sm text-gray-500">
                            Showing {data.rows.length} of {data.total_rows} events
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
                </>
            )}
        </div>
    )
}