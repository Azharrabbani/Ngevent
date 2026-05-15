import Pagination from "../../../../components/pagination";
import EventTable from "./eventTable";
import EventCard from "./eventCard";

interface Props {
    pendingEvent: boolean
    currentPage: number
    totalPage: number
    setCurrentPage: React.Dispatch<React.SetStateAction<number>>
}

export default function EventList({ pendingEvent, currentPage, totalPage, setCurrentPage }: Props) {


    return (
        <div className="bg-white rounded-3xl shadow-sm border border-gray-100 overflow-hidden">
            {/* Desktop Table */}
            <EventTable isPending={pendingEvent} />

            {/* Mobile Card View */}
            <EventCard isPending={pendingEvent} />

            <div className="flex flex-col md:flex-row items-center justify-between gap-4 px-6 md:px-8 border-t border-gray-100">
                <p className="text-sm text-gray-500">
                    Showing 1 to 4 of 12 events
                </p>

                <Pagination
                    currentPage={currentPage}
                    totalPage={totalPage}
                    onPrev={() => setCurrentPage((prev) => Math.max(prev - 1, 1))}
                    onNext={() => setCurrentPage((prev) => Math.min(prev + 1, totalPage))}
                    onCurrent={(page) => setCurrentPage(page)}
                />
            </div>
        </div>
    )
}