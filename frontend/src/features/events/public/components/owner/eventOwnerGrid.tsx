import { ShockIcon } from "../../../../../components/icon";
import type { OrganizerResponse } from "../../../../profile/types/profileResponse";
import EventOwnerCard from "./eventOwnerCards";


interface Props {
    organizers: OrganizerResponse[];
    onSelect?: (organizer: OrganizerResponse) => void;
}

export default function EventOwnerGrid({ organizers, onSelect }: Props) {
    if (organizers.length === 0) {
        return (
            <div className="mt-16 flex flex-col items-center justify-center text-center">
                <div className="w-16 h-16 rounded-2xl bg-slate-100 flex items-center justify-center mb-4">
                    <ShockIcon className="text-3xl text-indigo-500" />
                </div>
                <p className="text-slate-700 font-semibold">No event owners found</p>
                <p className="text-slate-400 text-sm mt-1">
                    Try a different search keyword
                </p>
            </div>
        );
    }

    return (
        <div className="mt-8 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
            {organizers.map((organizer) => (
                <EventOwnerCard
                    key={organizer.id}
                    organizer={organizer}
                    onClick={() => onSelect?.(organizer)}
                />
            ))}
        </div>
    );
}