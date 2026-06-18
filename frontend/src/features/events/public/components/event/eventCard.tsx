import { useNavigate } from "react-router-dom";
import type { EventsResponse } from "../../../types/eventResponse";
import { eventDateRange } from "../../../../../utils/dateConverter";
import { FiMapPin, FiCalendar } from "react-icons/fi";

interface Props {
    event: EventsResponse;
}

export default function EventCard({
    event,
}: Props) {
    const navigate = useNavigate();

    const eventDate = eventDateRange(
        event.start_date,
        event.end_date
    );

    return (
        <div
            onClick={() =>
                navigate(
                    `/events/${event.event.slug}`
                )
            }
            className="
                group
                bg-white rounded-2xl overflow-hidden
                border border-slate-200
                hover:border-blue-200
                hover:shadow-xl hover:shadow-blue-100/40
                transition-all duration-300
                cursor-pointer
            "
        >
            <div className="relative h-48 overflow-hidden bg-slate-100">
                <img
                    src={event.event.banner}
                    alt={event.event.name}
                    className="h-full w-full object-cover group-hover:scale-105 transition-transform duration-500"
                />
                <div className="absolute inset-0 bg-gradient-to-t from-black/30 to-transparent" />
            </div>

            <div className="p-4">
                <h3 className="font-semibold text-slate-900 text-sm leading-snug line-clamp-2 group-hover:text-blue-700 transition-colors">
                    {event.event.name}
                </h3>

                <div className="mt-3 space-y-1.5">
                    <div className="flex items-center gap-1.5 text-xs text-slate-500">
                        <FiCalendar className="w-3.5 h-3.5 shrink-0 text-slate-400" />
                        <span>{eventDate}</span>
                    </div>

                    <div className="flex items-center gap-1.5 text-xs text-slate-500">
                        <FiMapPin className="w-3.5 h-3.5 shrink-0 text-slate-400" />
                        <span className="truncate">{event.event_address.city}</span>
                    </div>
                </div>

                <div className="mt-3 pt-3 border-t border-slate-100 flex items-center justify-between">
                    <span className="text-xs text-slate-500 truncate">
                        by{" "}
                        <span className="font-medium text-slate-700">
                            {event.eo_profile.name}
                        </span>
                    </span>
                </div>
            </div>
        </div>
    );
}
