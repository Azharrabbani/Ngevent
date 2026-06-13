import { FiCalendar, FiMapPin } from "react-icons/fi";
import type { EventsResponse } from "../../../types/eventResponse";
import { toDateString } from "../../../../../utils/dateConverter";

interface Props {
    event: EventsResponse;
    onClick?: () => void;
}

export default function EventCard({ event, onClick }: Props) {
    const start = new Date(Number(event.start_time) * 1000);
    const date = toDateString(start, "long");
    return (
        <div
            onClick={onClick}
            className="
                bg-white rounded-2xl border border-slate-200 overflow-hidden
                hover:shadow-md hover:border-blue-200
                transition-all duration-200 cursor-pointer group
            "
        >
            <div className="relative h-44 bg-slate-100 overflow-hidden">
                {event.event?.banner ? (
                    <img
                        src={event.event?.banner}
                        alt={event.event?.name}
                        className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
                    />
                ) : (
                    <div className="w-full h-full bg-gradient-to-br from-slate-100 to-slate-200 flex items-center justify-center">
                        <FiCalendar className="w-10 h-10 text-slate-300" />
                    </div>
                )}
            </div>

            <div className="p-4">
                <h3 className="font-bold text-slate-900 text-base leading-snug group-hover:text-blue-600 transition line-clamp-2">
                    {event.event?.name}
                </h3>

                <div className="mt-3 space-y-1.5">
                    <div className="flex items-center gap-2 text-slate-500 text-xs">
                        <FiCalendar className="w-3.5 h-3.5 shrink-0" />
                        <span>{date}</span>
                    </div>
                    <div className="flex items-center gap-2 text-slate-500 text-xs">
                        <FiMapPin className="w-3.5 h-3.5 shrink-0" />
                        <span className="truncate">{event.event_address?.city || "—"}</span>
                    </div>
                </div>

                <div className="mt-4 pt-3 border-t border-slate-100 flex items-center gap-2">
                    {event.eo_profile?.photo_profile ? (
                        <img
                            src={event.eo_profile?.photo_profile}
                            alt={event.eo_profile?.name}
                            className="w-6 h-6 rounded-full object-cover border border-slate-200"
                        />
                    ) : (
                        <div className="w-6 h-6 rounded-full bg-slate-200 flex items-center justify-center text-slate-500 text-xs font-bold">
                            {event.eo_profile?.name?.charAt(0)?.toUpperCase() ?? "?"}
                        </div>
                    )}
                    <span className="text-xs text-slate-500">
                        <span className="font-medium text-slate-700">
                            {event.eo_profile?.name}
                        </span>
                    </span>
                </div>
            </div>
        </div>
    );
}