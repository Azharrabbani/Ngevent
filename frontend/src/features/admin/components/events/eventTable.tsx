import { useState } from "react";
import { IoCheckmark } from "react-icons/io5";
import { RxCross2 } from "react-icons/rx";
import { useNavigate } from "react-router-dom";
import type { PaginatedData } from "../../../../types/apiResponse";
import type { EventsResponse } from "../../../events/types/eventResponse";
import { FormatRelativeTime } from "../../../../utils/formatRelativeTime";
import { getStatusColor } from "../../../../utils/status";
import { eventDateRange } from "../../../../utils/dateConverter";
import { useReviewEvent } from "../../../events/hooks/useReviewEvent";
import ApproveConfirmModal from "./modal/approveConfirmModal";
import RejectReasonModal from "./modal/rejectReasonModal";
import { groupEventsByUrgency } from "../../utils/eventUrgency";
import UrgencyGroupHeader from "./header/urgencyGroupHeader";
import UrgencyBadge from "./badge/urgencyBadge";
import { useReviewUpdatedEvent } from "../../../events/hooks/useReviewUpdateEvent";

interface Props {
    status: string;
    data: PaginatedData<EventsResponse>;
    isReview: boolean;
    sort?: string;
    setSort?: React.Dispatch<React.SetStateAction<string | undefined>>;
    getUpdate?: boolean;
}

export default function EventTable({ status, isReview, data, getUpdate }: Props) {
    const navigate = useNavigate();
    const { mutateAsync: reviewEvent, isPending } = useReviewEvent();
    const { mutateAsync: reviewUpdatedEvent, isPending: isReviewUpdatePending } = useReviewUpdatedEvent();

    const [approveTarget, setApproveTarget] = useState<EventsResponse | null>(null);
    const [rejectTarget, setRejectTarget] = useState<EventsResponse | null>(null);

    const handleApprove = async () => {
        if (!approveTarget) return;
        if (getUpdate) {
            const updateEvent = approveTarget?.update_request_id || ""
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
            const updateEvent = rejectTarget?.update_request_id || ""
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
            <div className="hidden 2xl:block overflow-x-auto">
                <table className="w-full min-w-[800px]">
                    <thead className="bg-[#F8F9FC]">
                        <tr>
                            <th className="px-8 py-5 text-left text-sm font-semibold text-gray-600 border-b">Event Details</th>
                            <th className="px-8 py-5 text-left text-sm font-semibold text-gray-600 border-b">Event Owner</th>
                            <th className="px-8 py-5 text-left text-sm font-semibold text-gray-600 border-b">Submitted</th>
                            <th className={`px-8 py-5 text-sm font-semibold text-gray-600 border-b ${isReview && status === "pending" ? "text-center" : "text-left"}`}>
                                {isReview && status === "pending" ? "Actions" : "Status"}
                            </th>
                        </tr>
                    </thead>

                    <tbody>
                        {urgencyGroups
                            ? /* Grouped (pending review) */
                            urgencyGroups.map((group) => (
                                <>
                                    <UrgencyGroupHeader key={`hdr-${group.level}`} group={group} />

                                    {group.events.map((item) => {
                                        const eventDate = eventDateRange(
                                            item.start_date,
                                            item.end_date
                                        );
                                        const submitted = item.submitted_at ? FormatRelativeTime(item.submitted_at) : "-";

                                        return (
                                            <tr
                                                key={item.id}
                                                onClick={() => navigate(`/admin/events/review/${item.id}?status=${status}`)}
                                                className="hover:bg-gray-50 cursor-pointer transition-colors duration-200"
                                            >
                                                <td className="px-8 py-5 border-b border-gray-100">
                                                    <div className="flex items-center gap-4">
                                                        <div className="w-14 h-14 rounded-2xl bg-[#EEF0FF] overflow-hidden shrink-0">
                                                            <img src={item.event.banner} alt="event_banner" className="w-full h-full object-cover" />
                                                        </div>
                                                        <div>
                                                            <h1 className="font-semibold text-base text-gray-800">{item.event.name}</h1>
                                                            <p className="text-gray-500 text-sm">{eventDate}</p>
                                                            <div className="mt-1.5">
                                                                <UrgencyBadge
                                                                    urgency={item.urgency}
                                                                    badgeColor={group.badgeColor}
                                                                />
                                                            </div>
                                                        </div>
                                                    </div>
                                                </td>
                                                <td className="px-8 py-5 border-b border-gray-100">
                                                    <div>
                                                        <div className="flex items-center gap-1">
                                                            <h1 className="font-semibold text-gray-800">{item.eo_profile.name}</h1>
                                                            {item.eo_profile.status === "deactivated" && (
                                                                <span className="px-2 py-1 rounded-full text-xs bg-red-100 text-red-600">Deactivated</span>
                                                            )}
                                                        </div>
                                                        <p className="text-gray-500 text-sm">{item.eo_profile.email}</p>
                                                    </div>
                                                </td>
                                                <td className="px-8 py-5 border-b border-gray-100">
                                                    <span className="px-4 py-2 rounded-full text-sm bg-orange-100 text-orange-700">{submitted}</span>
                                                </td>
                                                <td className="px-8 py-5 border-b border-gray-100">
                                                    <div className="flex items-center justify-center gap-3">
                                                        <button
                                                            onClick={(e) => { e.stopPropagation(); navigate(`/admin/events/review/${item.id}?status=${status}`); }}
                                                            className="px-7 py-3 border border-gray-300 rounded-xl hover:bg-gray-100 transition-colors text-sm"
                                                        >
                                                            Review
                                                        </button>
                                                        <button
                                                            onClick={(e) => { e.stopPropagation(); setRejectTarget(item); }}
                                                            className="w-10 h-10 bg-white text-red-600 border border-red-600 rounded-lg hover:bg-red-600 hover:text-white transition-colors flex items-center justify-center"
                                                        >
                                                            <RxCross2 size={24} />
                                                        </button>
                                                        <button
                                                            onClick={(e) => { e.stopPropagation(); setApproveTarget(item); }}
                                                            className="w-10 h-10 bg-blue-500 text-white rounded-lg hover:bg-blue-700 transition-colors flex items-center justify-center"
                                                        >
                                                            <IoCheckmark size={24} />
                                                        </button>
                                                    </div>
                                                </td>
                                            </tr>
                                        );
                                    })}
                                </>
                            ))

                            : /* Non-pending */
                            data.rows.map((item) => {
                                const eventDate = eventDateRange(
                                    item.start_date,
                                    item.end_date
                                );
                                const submitted = item.submitted_at ? FormatRelativeTime(item.submitted_at) : "-";

                                return (
                                    <tr
                                        key={item.id}
                                        onClick={() => navigate(`/admin/events/review/${item.id}?status=${status}`)}
                                        className="hover:bg-gray-50 cursor-pointer transition-colors duration-200"
                                    >
                                        <td className="px-8 py-6 border-b border-gray-100">
                                            <div className="flex items-center gap-4">
                                                <div className="w-14 h-14 rounded-2xl bg-[#EEF0FF] overflow-hidden">
                                                    <img src={item.event.banner} alt="event_banner" className="w-full h-full object-cover" />
                                                </div>
                                                <div>
                                                    <h1 className="font-semibold text-lg text-gray-800">{item.event.name}</h1>
                                                    <p className="text-gray-500 text-sm">{eventDate}</p>
                                                </div>
                                            </div>
                                        </td>
                                        <td className="px-8 py-6 border-b border-gray-100">
                                            <div className="flex items-center gap-1">
                                                <h1 className="font-semibold text-gray-800">{item.eo_profile.name}</h1>
                                                {item.eo_profile.status === "deactivated" && (
                                                    <span className="px-2 py-1 rounded-full text-sm bg-red-100 text-red-600">Deactivated</span>
                                                )}
                                            </div>
                                            <p className="text-gray-500 text-sm">{item.eo_profile.email}</p>
                                        </td>
                                        <td className="px-8 py-6 border-b border-gray-100">
                                            <span className="px-4 py-2 rounded-full text-sm bg-orange-100 text-orange-700">{submitted}</span>
                                        </td>
                                        <td className="px-8 py-6 border-b border-gray-100">
                                            <span className={`font-medium ${getStatusColor(item?.event?.status ?? "")}`}>
                                                {item.event.status}
                                            </span>
                                        </td>
                                    </tr>
                                );
                            })
                        }
                    </tbody>
                </table>
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