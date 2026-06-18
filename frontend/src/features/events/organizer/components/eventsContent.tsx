import type { PaginatedData } from "../../../../types/apiResponse"
import type { EventsResponse } from "../../types/eventResponse"
import EventCard from "./eventCard";

interface Props {
    data: PaginatedData<EventsResponse> | undefined;
    loading: boolean;
};

export default function EventsContent({ data, loading }: Props) {
    return (
        <div className="p-12">
            <div className="space-y-2 text-center md:text-start">
                <h2 className="text-lg font-semibold">Your Events</h2>
                <p className="text-[#424654]">Showing <b>{data?.total_rows}</b> events</p>
            </div>

            <div className="mt-8 grid grid-cols-1 sm:grid-cols-2 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
                {loading ? (
                    <h1>Loading...</h1>
                ) : !data?.rows || data.rows.length === 0 ? (
                    <h1 className="text-gray-400 text-sm text-center md:text-start">
                        Events not found
                    </h1>
                ) : (
                    data.rows.map((data) => (
                        <EventCard
                            key={data.id}
                            id={data.id}
                            title={data.event.name}
                            startDate={data.start_date}
                            endDate={data.end_date}
                            startTime={data.start_time}
                            endTime={data.end_time}
                            location={data.event_address.city}
                            status={data.event.status}
                            image={data.event.banner}
                        />
                    ))
                )}
            </div>
        </div>
    )
}