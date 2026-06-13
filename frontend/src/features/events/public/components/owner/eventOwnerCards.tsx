import { useNavigate } from "react-router-dom";
import type { OrganizerResponse } from "../../../../profile/types/profileResponse";


interface Props {
    organizer: OrganizerResponse;
}

export default function EventOwnerCard({ organizer }: Props) {
    const navigate = useNavigate();

    return (
        <div
            onClick={() => navigate(`/event-owner/${organizer.slug}`)}
            className="
                bg-white rounded-2xl border border-slate-200
                overflow-hidden
                hover:shadow-md hover:border-blue-200
                transition-all duration-200
                cursor-pointer
                group
                flex
            "
        >
            <div className="w-24 shrink-0">
                {organizer.photo_profile ? (
                    <img
                        src={organizer.photo_profile}
                        alt={organizer.name}
                        className="
                            w-full h-full
                            object-cover
                        "
                    />
                ) : (
                    <div
                        className="
                            w-full h-full min-h-[120px]
                            bg-gradient-to-br from-blue-50 to-indigo-100
                            flex items-center justify-center
                            text-blue-600 font-bold text-4xl
                        "
                    >
                        {organizer.name?.charAt(0)?.toUpperCase() ?? "?"}
                    </div>
                )}
            </div>

            <div className="flex-1 p-4 flex flex-col justify-center min-w-0">
                <h3 className="text-sm font-semibold text-slate-900 truncate group-hover:text-blue-600 transition">
                    {organizer.name}
                </h3>

                <p className="text-xs text-slate-500 truncate mt-1">
                    {organizer.email}
                </p>

                <p className="text-xs font-medium text-blue-500 mt-3">
                    {organizer.event_count ?? 0} events
                </p>
            </div>
        </div>
    );
}