import { useNavigate, useParams } from "react-router-dom";
import { GoArrowLeft } from "react-icons/go";
import { GoDotFill } from "react-icons/go";
import { FiClock, FiCalendar, FiMapPin, FiTag, FiAlertCircle, FiInfo } from "react-icons/fi";
import { MdOutlineLocationOn } from "react-icons/md";
import { FaBullhorn } from "react-icons/fa";
import { useState } from "react";
import { useGetEventByID } from "../../hooks/useGetEventByID";
import { useCancelEvent } from "../../hooks/useCancelEvent";
import { useDeleteEvent } from "../../hooks/useDeleteEvent";
import { timeRange } from "../../../../utils/timeRange";
import Button from "../../../../components/Button";
import { eventDateRange } from "../../../../utils/dateConverter";

const statusColorMap: Record<string, { text: string; bg: string; dot: string }> = {
    active: { text: "text-[#0040A1]", bg: "bg-blue-50", dot: "text-[#0040A1]" },
    pending: { text: "text-amber-600", bg: "bg-amber-50", dot: "text-amber-500" },
    rejected: { text: "text-red-600", bg: "bg-red-50", dot: "text-red-500" },
    cancelled: { text: "text-red-600", bg: "bg-red-50", dot: "text-red-500" },
    done: { text: "text-green-600", bg: "bg-green-50", dot: "text-green-500" },
    draft: { text: "text-gray-500", bg: "bg-gray-100", dot: "text-gray-400" },
};

export default function EventViewPage() {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();

    const { data: eventData, isLoading } = useGetEventByID(id ?? "");
    const cancelEventMutation = useCancelEvent();
    const deleteEventMutation = useDeleteEvent();

    const [confirmAction, setConfirmAction] = useState<"cancel" | "delete" | null>(null);

    const handleCancel = async () => {
        if (!id) return;
        await cancelEventMutation.mutateAsync(id);
        navigate(-1);
    };

    const handleDelete = async () => {
        if (!id) return;
        await deleteEventMutation.mutateAsync(id);
        navigate(-1);
    };

    if (isLoading) {
        return (
            <div className="min-h-screen bg-[#F4F7FB] flex items-center justify-center">
                <div className="flex flex-col items-center gap-3">
                    <div className="w-10 h-10 border-4 border-[#003B95] border-t-transparent rounded-full animate-spin" />
                    <p className="text-gray-500 text-sm">Loading event...</p>
                </div>
            </div>
        );
    }

    if (!eventData) {
        return (
            <div className="min-h-screen bg-[#F4F7FB] flex items-center justify-center">
                <div className="text-center space-y-3">
                    <p className="text-gray-500 text-lg">Event not found.</p>
                    <button
                        onClick={() => navigate(-1)}
                        className="text-[#003B95] underline text-sm"
                    >
                        Go back
                    </button>
                </div>
            </div>
        );
    }

    const event = eventData.event;
    const address = eventData.event_address;
    const status = event.status;

    const start = new Date(Number(eventData.start_time) * 1000);
    const end = new Date(Number(eventData.end_time) * 1000);

    const eventDate = eventDateRange(
        eventData.start_date,
        eventData.end_date
    );

    const timeRangeVal = timeRange(start, end);
    const statusStyle = statusColorMap[status] ?? statusColorMap["draft"];

    const renderButtons = () => {
        if (status === "active") {
            return (
                <div className="flex flex-col sm:flex-row gap-3 w-full sm:w-auto">
                    <Button
                        type="button"
                        onClick={() => setConfirmAction("cancel")}
                        disabled={cancelEventMutation.isPending}
                        className="w-full sm:w-auto rounded-xl px-8 py-3 text-white font-semibold bg-red-500 hover:bg-red-600 transition"
                    >
                        {cancelEventMutation.isPending ? "Canceling..." : "Cancel Event"}
                    </Button>
                    <Button
                        type="button"
                        onClick={() => navigate(`/organizer/event/edit/${id}`)}
                        disabled={event.request_updates}
                        className="w-full sm:w-auto rounded-xl px-8 py-3 text-white font-semibold bg-[#003B95] hover:bg-[#004ec2] transition disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        Edit Event
                    </Button>
                </div>
            );
        }

        if (status === "draft") {
            return (
                <div className="flex flex-col sm:flex-row gap-3 w-full sm:w-auto">
                    <Button
                        type="button"
                        onClick={() => setConfirmAction("delete")}
                        disabled={deleteEventMutation.isPending}
                        className="w-full sm:w-auto rounded-xl px-8 py-3 text-white font-semibold bg-red-500 hover:bg-red-600 transition"
                    >
                        {deleteEventMutation.isPending ? "Deleting..." : "Delete Event"}
                    </Button>
                    <Button
                        type="button"
                        onClick={() => navigate(`/organizer/event/edit/${id}`)}
                        className="w-full sm:w-auto rounded-xl px-8 py-3 text-white font-semibold bg-[#003B95] hover:bg-[#004ec2] transition"
                    >
                        Edit Event
                    </Button>
                </div>
            );
        }

        return null;
    };

    return (
        <div className="min-h-screen bg-[#F4F7FB]">
            <div className="w-full max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12 space-y-8">

                {/* Page Header */}
                <div className="flex flex-col sm:flex-row gap-4 sm:gap-6 sm:items-start">
                    <FaBullhorn className="hidden sm:block text-[#003B95] mt-1 shrink-0" size={40} />
                    <div className="flex-1 min-w-0">
                        <span className="flex items-center gap-2">
                            <GoArrowLeft
                                className="cursor-pointer hover:-translate-x-1 transition duration-150 shrink-0"
                                onClick={() => navigate(-1)}
                                size={26}
                            />
                            <h1 className="text-2xl sm:text-3xl lg:text-4xl text-[#1E293B] font-bold truncate">
                                Event Details
                            </h1>
                        </span>
                        <p className="text-gray-500 text-sm sm:text-base mt-1 ml-8 sm:ml-0">
                            View the full details of this event
                        </p>
                    </div>
                </div>

                <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
                    <div className="relative w-full h-52 sm:h-72 lg:h-96 bg-gray-100">
                        <img
                            src={
                                event.banner ||
                                "https://t4.ftcdn.net/jpg/16/79/44/21/360_F_1679442196_OEsi0AFKie6hYMBpvmXwwRgRYGV4U6Lz.jpg"
                            }
                            alt={event.name}
                            className="w-full h-full object-cover"
                        />
                        <div
                            className={`absolute top-4 left-4 flex items-center gap-1.5 px-3 py-1.5 rounded-full text-sm font-medium backdrop-blur-sm bg-white/90 shadow-sm ${statusStyle.text}`}
                        >
                            <GoDotFill size={12} className={statusStyle.dot} />
                            <span className="capitalize">{status}</span>
                        </div>
                    </div>

                    <div className="p-5 sm:p-8 lg:p-10 space-y-8">
                        <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4">
                            <h2 className="text-2xl sm:text-3xl font-bold text-[#1E293B] leading-tight">
                                {event.name}
                            </h2>
                            {renderButtons()}
                        </div>

                        {event.request_updates && (
                            <div className="flex gap-3 p-4 rounded-xl bg-amber-50 border border-amber-200">
                                <div className="shrink-0 mt-0.5">
                                    <FiInfo className="text-amber-500" size={18} />
                                </div>
                                <div>
                                    <p className="text-sm font-semibold text-amber-600 mb-0.5">
                                        Update Under Review
                                    </p>
                                    <p className="text-sm text-amber-500 leading-relaxed">
                                        You have a pending update request for this event. Editing is
                                        disabled until the current update has been reviewed by our
                                        admins.
                                    </p>
                                </div>
                            </div>
                        )}

                        {/* ── Rejection reason ─────────────────────────────── */}
                        {status === "rejected" && event.rejected_reason && (
                            <div className="flex gap-3 p-4 rounded-xl bg-red-50 border border-red-200">
                                <div className="shrink-0 mt-0.5">
                                    <FiAlertCircle className="text-red-500" size={18} />
                                </div>
                                <div>
                                    <p className="text-sm font-semibold text-red-600 mb-0.5">
                                        Rejection Reason
                                    </p>
                                    <p className="text-sm text-red-500 leading-relaxed">
                                        {event.rejected_reason}
                                    </p>
                                </div>
                            </div>
                        )}

                        <div className="grid grid-cols-1 sm:grid-cols-2 gap-5">
                            <div className="flex items-start gap-3 p-4 rounded-xl bg-[#F4F7FB]">
                                <div className="p-2 bg-[#003B95]/10 rounded-lg shrink-0">
                                    <FiCalendar className="text-[#003B95]" size={18} />
                                </div>
                                <div>
                                    <p className="text-xs text-gray-400 font-medium uppercase tracking-wide mb-0.5">Date</p>
                                    <p className="text-sm sm:text-base font-semibold text-[#1E293B]">{eventDate}</p>
                                </div>
                            </div>

                            <div className="flex items-start gap-3 p-4 rounded-xl bg-[#F4F7FB]">
                                <div className="p-2 bg-[#003B95]/10 rounded-lg shrink-0">
                                    <FiClock className="text-[#003B95]" size={18} />
                                </div>
                                <div>
                                    <p className="text-xs text-gray-400 font-medium uppercase tracking-wide mb-0.5">Time</p>
                                    <p className="text-sm sm:text-base font-semibold text-[#1E293B]">{timeRangeVal} WIB</p>
                                </div>
                            </div>

                            <div className="flex items-start gap-3 p-4 rounded-xl bg-[#F4F7FB]">
                                <div className="p-2 bg-[#003B95]/10 rounded-lg shrink-0">
                                    <MdOutlineLocationOn className="text-[#003B95]" size={18} />
                                </div>
                                <div>
                                    <p className="text-xs text-gray-400 font-medium uppercase tracking-wide mb-0.5">City</p>
                                    <p className="text-sm sm:text-base font-semibold text-[#1E293B]">{address.city}</p>
                                </div>
                            </div>

                            <div className="flex items-start gap-3 p-4 rounded-xl bg-[#F4F7FB]">
                                <div className="p-2 bg-[#003B95]/10 rounded-lg shrink-0">
                                    <FiTag className="text-[#003B95]" size={18} />
                                </div>
                                <div>
                                    <p className="text-xs text-gray-400 font-medium uppercase tracking-wide mb-1">Categories</p>
                                    <div className="flex flex-wrap gap-1.5">
                                        {event.categories?.length ? (
                                            event.categories.map((cat) => (
                                                <span
                                                    key={cat.id}
                                                    className="px-2.5 py-0.5 bg-[#003B95]/10 text-[#003B95] text-xs font-medium rounded-full"
                                                >
                                                    {cat.name}
                                                </span>
                                            ))
                                        ) : (
                                            <span className="text-sm text-gray-400">No categories</span>
                                        )}
                                    </div>
                                </div>
                            </div>
                        </div>

                        <div className="space-y-2">
                            <div className="flex items-center gap-2">
                                <FiMapPin className="text-[#003B95]" size={16} />
                                <h3 className="text-sm font-semibold text-[#003B95] uppercase tracking-wide">
                                    Full Address
                                </h3>
                            </div>
                            <div className="bg-[#F4F7FB] rounded-xl p-4 text-sm text-[#424654] leading-relaxed">
                                <p className="font-medium text-[#1E293B]">{address.address}</p>
                                {address.detail_address && (
                                    <p className="mt-1 text-gray-500">{address.detail_address}</p>
                                )}
                            </div>
                        </div>

                        {event.description && (
                            <div className="space-y-2">
                                <h3 className="text-sm font-semibold text-[#003B95] uppercase tracking-wide">
                                    Description
                                </h3>
                                <div
                                    className="prose prose-sm max-w-none text-[#424654] bg-[#F4F7FB] rounded-xl p-4 sm:p-6"
                                    dangerouslySetInnerHTML={{ __html: event.description }}
                                />
                            </div>
                        )}
                    </div>
                </div>
            </div>

            {confirmAction && (
                <div
                    className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4"
                    onClick={() => setConfirmAction(null)}
                >
                    <div
                        className="bg-white rounded-2xl p-6 max-w-sm w-full shadow-xl"
                        onClick={(e) => e.stopPropagation()}
                    >
                        <h2 className="text-lg font-bold text-gray-800 mb-2">
                            {confirmAction === "cancel" ? "Cancel Event?" : "Delete Event?"}
                        </h2>
                        <p className="text-sm text-gray-500 mb-6">
                            {confirmAction === "cancel"
                                ? "This action cannot be undone. The event will be canceled."
                                : "This will permanently delete the event draft. This action cannot be undone."}
                        </p>
                        <div className="flex gap-3 justify-end">
                            <button
                                onClick={() => setConfirmAction(null)}
                                className="px-4 py-2 rounded-lg border border-gray-300 hover:bg-gray-100 text-sm transition"
                            >
                                Go Back
                            </button>
                            <button
                                onClick={() => {
                                    setConfirmAction(null);
                                    confirmAction === "cancel" ? handleCancel() : handleDelete();
                                }}
                                className="px-4 py-2 rounded-lg bg-red-500 hover:bg-red-600 text-white text-sm transition"
                            >
                                {confirmAction === "cancel" ? "Yes, Cancel" : "Yes, Delete"}
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}