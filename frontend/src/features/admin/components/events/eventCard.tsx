import { useState } from "react";
import { IoCheckmark } from "react-icons/io5";
import { RxCross2 } from "react-icons/rx";
import { useNavigate } from "react-router-dom";
import type { EventsResponse } from "../../../events/types/eventResponse";
import type { PaginatedData } from "../../../../types/apiResponse";
import { FormatRelativeTime } from "../../../../utils/formatRelativeTime";
import { getStatusColor } from "../../../../utils/status";
import { eventDateRange, toDateString } from "../../../../utils/dateConverter";
import { useReviewEvent } from "../../../events/hooks/useReviewEvent";
import ApproveConfirmModal from "./modal/approveConfirmModal";
import RejectReasonModal from "./modal/rejectReasonModal";
import { groupEventsByUrgency, type UrgencyGroup } from "../../utils/eventUrgency";
import UrgencyBadge from "./badge/urgencyBadge";
import { CircleIcon } from "../../../../components/icon";
import { useReviewUpdatedEvent } from "../../../events/hooks/useReviewUpdateEvent";

interface Props {
    status: string;
    data: PaginatedData<EventsResponse>;
    isReview: boolean;
    sort?: string;
    setSort?: React.Dispatch<React.SetStateAction<string | undefined>>;
    getUpdate?: boolean;
}

const groupStripe: Record<string, string> = {
    critical: "border-l-red-500",
    urgent: "border-l-orange-400",
    warning: "border-l-yellow-400",
    normal: "border-l-green-400",
};

interface CardGroupSectionProps {
    group: UrgencyGroup;
    isReview: boolean;
    status: string;
    onNavigate: (id: string) => void;
    onApprove: (item: EventsResponse) => void;
    onReject: (item: EventsResponse) => void;
}

function CardGroupSection({ group, isReview, status, onNavigate, onApprove, onReject }: CardGroupSectionProps) {
    return (
        <div className="space-y-3">
            <div className="flex items-center gap-2 px-1">
                <CircleIcon className={`w-2 h-2 shrink-0 ${group.iconColor}`} size={8} />
                <span className={`text-sm font-bold ${group.color}`}>
                    {group.title}
                </span>
                <span className="text-xs text-gray-400">({group.events.length})</span>
            </div>

            {group.events.map((item) => {
                const eventDate = eventDateRange(
                    item.start_date,
                    item.end_date
                );
                const submitted = item.submitted_at ? FormatRelativeTime(item.submitted_at) : "-";

                return (
                    <div
                        key={item.id}
                        onClick={() => onNavigate(item.id)}
                        className={`
                            w-full cursor-pointer transition-colors duration-200
                            border border-gray-100 border-l-4 rounded-2xl p-4 shadow-sm bg-white
                            hover:bg-gray-50
                            ${groupStripe[group.level]}
                        `}
                    >
                        <div className="flex items-start gap-4">
                            <div className="w-14 h-14 shrink-0">
                                <img src={item.event.banner} alt="event_banner" className="w-full h-full object-cover rounded-2xl" />
                            </div>
                            <div className="flex-1 min-w-0">
                                <div className="flex items-start justify-between gap-3">
                                    <h1 className="font-semibold text-base text-gray-800 leading-snug">{item.event.name}</h1>
                                    {(!isReview || status !== "pending") && (
                                        <span className={`text-sm font-medium whitespace-nowrap ${getStatusColor(item.event.status ?? "")}`}>
                                            {item.event.status}
                                        </span>
                                    )}
                                </div>
                                <p className="text-sm text-gray-500 mt-0.5">{eventDate}</p>

                                <div className="mt-1.5">
                                    <UrgencyBadge urgency={item.urgency} badgeColor={group.badgeColor} />
                                </div>

                                <div className="mt-3">
                                    <p className="text-xs text-gray-400">Organizer</p>
                                    <div className="flex items-center gap-1">
                                        <h2 className="text-sm font-medium text-gray-700">{item.eo_profile.name}</h2>
                                        {item.eo_profile.status === "deactivated" && (
                                            <span className="px-2 py-1 rounded-full text-xs bg-red-100 text-red-600">Deactivated</span>
                                        )}
                                    </div>
                                    <p className="text-sm text-gray-500 break-all">{item.eo_profile.email}</p>
                                </div>

                                <div className="mt-4 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
                                    <span className="w-fit px-3 py-1 rounded-full text-xs bg-orange-100 text-orange-700">{submitted}</span>
                                    {isReview && status === "pending" && (
                                        <div className="flex items-center gap-2">
                                            <button
                                                onClick={(e) => { e.stopPropagation(); onNavigate(item.id); }}
                                                className="flex-1 sm:flex-none px-5 py-2 border border-gray-300 rounded-xl text-sm hover:bg-gray-100 transition-colors"
                                            >
                                                Review
                                            </button>
                                            <button
                                                onClick={(e) => { e.stopPropagation(); onReject(item); }}
                                                className="w-10 h-10 bg-white text-red-600 border border-red-600 rounded-lg hover:bg-red-600 hover:text-white transition-colors flex items-center justify-center"
                                            >
                                                <RxCross2 size={20} />
                                            </button>
                                            <button
                                                onClick={(e) => { e.stopPropagation(); onApprove(item); }}
                                                className="w-10 h-10 bg-blue-500 text-white rounded-lg hover:bg-blue-700 transition-colors flex items-center justify-center"
                                            >
                                                <IoCheckmark size={20} />
                                            </button>
                                        </div>
                                    )}
                                </div>
                            </div>
                        </div>
                    </div>
                );
            })}
        </div>
    );
}

export default function EventCard({ status, isReview, data, getUpdate }: Props) {
    const navigate = useNavigate();
    const { mutateAsync: reviewEvent, isPending } = useReviewEvent();
    const { mutateAsync: reviewUpdatedEvent, isPending: isReviewUpdatePending } = useReviewUpdatedEvent();

    const [approveTarget, setApproveTarget] = useState<EventsResponse | null>(null);
    const [rejectTarget, setRejectTarget] = useState<EventsResponse | null>(null);

    const handleApprove = async () => {
        if (!approveTarget) return;
        if (getUpdate) {
            const updateEvent = approveTarget?.update_request_id || "";
            await reviewUpdatedEvent({ id: updateEvent, status: "approved" });
        } else {
            await reviewEvent({ id: approveTarget.id, status: "active" });
        }
        navigate("/admin/events/pending");
        setApproveTarget(null);
    };

    const handleReject = async (reason: string) => {
        if (!rejectTarget) return;
        if (getUpdate) {
            const updateEvent = rejectTarget?.update_request_id || "";
            await reviewUpdatedEvent({ id: updateEvent, status: "rejected", reason });
        } else {
            await reviewEvent({ id: rejectTarget.id, status: "rejected", reason });
        }
        navigate("/admin/events/pending");
        setRejectTarget(null);
    };

    // Only grouped for pending events
    const urgencyGroups = isReview && status === "pending"
        ? groupEventsByUrgency(data.rows)
        : null;

    return (
        <>
            <div className="2xl:hidden w-full p-4 space-y-6">
                {urgencyGroups
                    ? urgencyGroups.map((group) => (
                        <CardGroupSection
                            key={group.level}
                            group={group}
                            isReview={isReview}
                            status={status}
                            onNavigate={(id) => navigate(`/admin/events/review/${id}?status=${status}`)}
                            onApprove={setApproveTarget}
                            onReject={setRejectTarget}
                        />
                    ))
                    // Non pending
                    : data.rows.map((item) => {
                        const start = new Date(Number(item.start_time) * 1000);
                        const date = toDateString(start, "short");
                        const submitted = item.submitted_at ? FormatRelativeTime(item.submitted_at) : "-";

                        return (
                            <div
                                key={item.id}
                                onClick={() => navigate(`/admin/events/review/${item.id}?status=${status}`)}
                                className="w-full hover:bg-gray-50 cursor-pointer transition-colors duration-200 border border-gray-100 rounded-2xl p-4 shadow-sm bg-white"
                            >
                                <div className="flex items-start gap-4">
                                    <div className="w-14 h-14 shrink-0">
                                        <img src={item.event.banner} alt="event_banner" className="w-full h-full object-cover rounded-2xl" />
                                    </div>
                                    <div className="flex-1 min-w-0">
                                        <div className="flex items-start justify-between gap-3">
                                            <h1 className="font-semibold text-base text-gray-800 leading-snug">{item.event.name}</h1>
                                            <span className={`text-sm font-medium whitespace-nowrap ${getStatusColor(item.event.status ?? "")}`}>
                                                {item.event.status}
                                            </span>
                                        </div>
                                        <p className="text-sm text-gray-500 mt-1">{date}</p>
                                        <div className="mt-3">
                                            <p className="text-xs text-gray-400">Organizer</p>
                                            <h2 className="text-sm font-medium text-gray-700">{item.eo_profile.name}</h2>
                                            <p className="text-sm text-gray-500 break-all">{item.eo_profile.email}</p>
                                        </div>
                                        <div className="mt-4">
                                            <span className="w-fit px-3 py-1 rounded-full text-xs bg-orange-100 text-orange-700">{submitted}</span>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        );
                    })
                }
            </div>

            <ApproveConfirmModal
                isOpen={!!approveTarget}
                eventName={approveTarget?.event?.name ?? ""}
                onClose={() => setApproveTarget(null)}
                onView={() => {
                    if (approveTarget) navigate(`/admin/events/review/${approveTarget.id}?status=${status}`);
                    setApproveTarget(null);
                }}
                onApprove={handleApprove}
                isPending={isPending || isReviewUpdatePending}
            />

            <RejectReasonModal
                isOpen={!!rejectTarget}
                onClose={() => setRejectTarget(null)}
                onSubmit={handleReject}
                isPending={isPending || isReviewUpdatePending}
            />
        </>
    );
}