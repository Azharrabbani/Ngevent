import { NoEventsIcon } from "../../../../../components/icon";
import type { EventsResponse } from "../../../types/eventResponse";
import EventCard from "./eventCard";

interface Props {
    events: EventsResponse[] | undefined;
}

export default function EventGrid({ events }: Props) {
    if (!events || !events.length) {
        return (
            <div className="py-24 flex flex-col items-center text-center">
                <div className="w-16 h-16 rounded-full bg-slate-100 flex items-center justify-center mb-4">
                    <NoEventsIcon className="text-3xl text-indigo-500" />
                </div>
                <h3 className="text-lg font-semibold text-slate-800">
                    Event not found
                </h3>
                <p className="mt-1 text-sm text-slate-500 max-w-xs">
                    Try changing the category, location, or date filter to find other events
                </p>
            </div>
        );
    }

    return (
        <div
            className="mt-8 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5"
        >
            {events.map((event) => (
                <EventCard
                    key={event.id}
                    event={event}
                />
            ))}
        </div>
    );
}
