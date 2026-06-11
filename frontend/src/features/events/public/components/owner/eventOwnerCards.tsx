import type { OrganizerResponse } from "../../../../profile/types/profileResponse";


interface Props {
    organizer: OrganizerResponse;
    onClick?: () => void;
}

export default function EventOwnerCard({ organizer, onClick }: Props) {
    return (
        <div
            onClick={onClick}
            className="
                bg-white rounded-2xl border border-slate-200
                p-4 flex items-center gap-4
                hover:shadow-md hover:border-blue-200
                transition-all duration-200
                cursor-pointer
                group
            "
        >
            <div className="shrink-0">
                {organizer.photo_profile ? (
                    <img
                        src={organizer.photo_profile}
                        alt={organizer.name}
                        className="
                            w-16 h-16 rounded-xl object-cover
                            border border-slate-100
                            group-hover:border-blue-200 transition
                        "
                    />
                ) : (
                    <div
                        className="
                            w-16 h-16 rounded-xl
                            bg-gradient-to-br from-blue-50 to-indigo-100
                            border border-slate-100
                            flex items-center justify-center
                            text-blue-600 font-bold text-xl
                        "
                    >
                        {organizer.name?.charAt(0)?.toUpperCase() ?? "?"}
                    </div>
                )}
            </div>

            <div className="flex-1 min-w-0">
                <h3 className="text-sm font-semibold text-slate-900 truncate group-hover:text-blue-600 transition">
                    {organizer.name}
                </h3>
                <p className="text-xs text-slate-500 truncate mt-0.5">
                    {organizer.email}
                </p>
                <p className="text-xs font-medium text-blue-500 mt-2">
                    {organizer.event_count ?? 0} event
                </p>
            </div>
        </div>
    );
}