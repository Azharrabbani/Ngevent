import { useState } from "react";
import { IoCheckmark } from "react-icons/io5";
import { RxCross2 } from "react-icons/rx";
import { useNavigate } from "react-router-dom";
import type { EventsResponse } from "../../../events/types/eventResponse";
import type { PaginatedData } from "../../../../types/apiResponse";
import { FormatRelativeTime } from "../../../../utils/formatRelativeTime";
import { getStatusColor } from "../../../../utils/status";
import { toDateString } from "../../../../utils/dateConverter";
import { useReviewEvent } from "../../../events/hooks/useReviewEvent";
import ApproveConfirmModal from "./modal/approveConfirmModal";
import RejectReasonModal from "./modal/rejectReasonModal";

interface Props {
    status: string;
    data: PaginatedData<EventsResponse> | undefined;
    isLoading: boolean;
    isReview: boolean;
    sort?: string;
    setSort?: React.Dispatch<React.SetStateAction<string | undefined>>;
}

export default function EventCard({ status, isReview, data, isLoading }: Props) {
    const navigate = useNavigate();
    const { mutateAsync: reviewEvent, isPending } = useReviewEvent();

    const [approveTarget, setApproveTarget] = useState<EventsResponse | null>(null);
    const [rejectTarget, setRejectTarget] = useState<EventsResponse | null>(null);

    const handleApprove = async () => {
        if (!approveTarget) return;
        await reviewEvent({ id: approveTarget.id, status: "active" });
        setApproveTarget(null);
    };

    const handleReject = async (reason: string) => {
        if (!rejectTarget) return;
        await reviewEvent({ id: rejectTarget.id, status: "rejected", reason });
        setRejectTarget(null);
    };

    return (
        <>
            <div className="xl:hidden p-4 space-y-4 max-h-150 sm:max-h-125 overflow-y-auto">
                {isLoading ? (
                    <h1>Loading...</h1>
                ) : !data?.rows || data.rows.length === 0 ? (
                    <div className="p-8">
                        <h1 className="text-gray-400 text-sm text-center md:text-start">Events not found</h1>
                    </div>
                ) : (
                    <>
                        {data.rows.map((item) => {
                            const start = new Date(Number(item?.start_time) * 1000);
                            const date = toDateString(start, "short");
                            const submitted = FormatRelativeTime(item?.created_at);

                            return (
                                <div
                                    key={item.id}
                                    onClick={() => navigate(`/admin/events/review/${item.id}?status=${status}`)}
                                    className="hover:bg-gray-50 cursor-pointer transition-colors duration-200 border border-gray-100 rounded-2xl p-4 shadow-sm bg-white"
                                >
                                    <div className="flex items-start gap-4">
                                        <div className="w-14 h-14 rounded-2xl bg-[#EEF0FF] flex items-center justify-center shrink-0">
                                            <img src={item.event.banner} alt="event_banner" className="w-full h-full object-cover" />
                                        </div>
                                        <div className="flex-1 min-w-0">
                                            <div className="flex items-start justify-between gap-3">
                                                <h1 className="font-semibold text-base text-gray-800 leading-snug">{item?.event?.name}</h1>
                                                {(!isReview || status !== "pending") && (
                                                    <span className={`text-sm font-medium whitespace-nowrap ${getStatusColor(item?.event?.status ?? "")}`}>
                                                        {item?.event?.status}
                                                    </span>
                                                )}
                                            </div>
                                            <p className="text-sm text-gray-500 mt-1">{date}</p>
                                            <div className="mt-4">
                                                <p className="text-xs text-gray-400">Organizer</p>
                                                <h2 className="text-sm font-medium text-gray-700">{item?.eo_profile?.name}</h2>
                                                <p className="text-sm text-gray-500 break-all">{item?.eo_profile?.email}</p>
                                            </div>
                                            <div className="mt-5 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
                                                <span className="w-fit px-3 py-1 rounded-full text-xs bg-orange-100 text-orange-700">{submitted}</span>
                                                {isReview && status === "pending" && (
                                                    <div className="flex items-center gap-2">
                                                        <button
                                                            onClick={(e) => { e.stopPropagation(); navigate(`/admin/events/review/${item.id}?status=${status}`); }}
                                                            className="flex-1 sm:flex-none px-5 py-2 border border-gray-300 rounded-xl text-sm hover:bg-gray-100 transition-colors"
                                                        >
                                                            Review
                                                        </button>
                                                        <button
                                                            onClick={(e) => { e.stopPropagation(); setRejectTarget(item); }}
                                                            className="w-10 h-10 bg-white text-red-600 border border-red-600 rounded-lg hover:bg-red-600 hover:text-white transition-colors flex items-center justify-center"
                                                        >
                                                            <RxCross2 size={20} />
                                                        </button>
                                                        <button
                                                            onClick={(e) => { e.stopPropagation(); setApproveTarget(item); }}
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
                    </>
                )}
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
                isPending={isPending}
            />

            <RejectReasonModal
                isOpen={!!rejectTarget}
                onClose={() => setRejectTarget(null)}
                onSubmit={handleReject}
                isPending={isPending}
            />
        </>
    );
}