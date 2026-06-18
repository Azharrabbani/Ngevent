import { EventIcon, PinIcon } from "../../../../../components/icon";
import { eventDateRange } from "../../../../../utils/dateConverter";

interface Props {
    name: string;
    startDate: number;
    endDate: number;
    startTime: number;
    endTime: number;
    detailAddress: string;
    city: string;
    distance?: string;
}

function formatEventTime(unix: number): string {
    const date = new Date(unix * 1000);
    return date.toLocaleTimeString("en-US", {
        hour: "2-digit",
        minute: "2-digit",
        hour12: false,
    });
}

function getTimezone(): string {
    return Intl.DateTimeFormat().resolvedOptions().timeZone
        .split("/")
        .pop()
        ?.toUpperCase() ?? "LOCAL";
}

export default function EventInfo({
    name,
    startDate,
    endDate,
    startTime,
    endTime,
    detailAddress,
    city,
    distance,
}: Props) {
    const startStr = formatEventTime(startTime);
    const endStr = formatEventTime(endTime);
    const tz = getTimezone();
    const eventDate = eventDateRange(
        startDate,
        endDate
    );

    return (
        <div className="bg-white rounded-2xl border border-slate-200 p-6 space-y-5">
            <div>
                <h1 className="text-2xl font-bold text-slate-900 leading-tight">{name}</h1>
            </div>

            <div className="space-y-3">
                <div className="flex items-start gap-3">
                    <div className="mt-0.5 w-8 h-8 rounded-lg bg-blue-50 flex items-center justify-center shrink-0">
                        <EventIcon className="w-4 h-4 text-blue-600" />
                    </div>
                    <div>
                        <p className="text-sm font-semibold text-slate-800">{eventDate}</p>
                        <p className="text-xs text-slate-500 mt-0.5">
                            {startStr} – {endStr}{" WIB "}
                            <span className="font-medium text-slate-600">{tz}</span>
                        </p>
                    </div>
                </div>

                <div className="flex items-start gap-3">
                    <div className="mt-0.5 w-8 h-8 rounded-lg bg-blue-50 flex items-center justify-center shrink-0">
                        <PinIcon className="w-4 h-4 text-blue-600" />
                    </div>
                    <div>
                        <p className="text-sm font-semibold text-slate-800">{detailAddress}</p>
                        <p className="text-xs text-slate-500 mt-0.5">
                            {city}
                            {distance && (
                                <span className="ml-1 text-blue-600 font-medium">
                                    ({distance} away)
                                </span>
                            )}
                        </p>
                    </div>
                </div>
            </div>
        </div>
    );
}