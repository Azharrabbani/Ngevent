import { FaRegUser } from "react-icons/fa";
import { RiCalendarEventFill, RiCalendarCheckFill, RiTimeFill } from "react-icons/ri";
import { MdCancel } from "react-icons/md";
import { useListUsers } from "../../profile/hooks/useListUsers";
import AdminSidebar from "../components/sideBar";
import { useNavigate } from "react-router-dom";
import { useGetEvents } from "../../events/hooks/useGetEvents";
import { useGetAllUpdatedEvents } from "../../events/hooks/useGetAllUpdatedEvent";

export default function AdminDashboard() {
    const organizers = useListUsers({
        role: "event organizer",
        pagination: { page: 1, limit: 1 },
    });

    const admins = useListUsers({
        role: "admin",
        pagination: { page: 1, limit: 1 },
    });

    const { data: activeEvents } = useGetEvents({
        status: "active",
        pagination: { page: 1, limit: 1 },
    });

    const { data: pendingUpdates } = useGetAllUpdatedEvents({
        status: "pending",
        pagination: { page: 1, limit: 1 },
    });

    const { data: pendingEvents } = useGetEvents({
        status: "pending",
        pagination: { page: 1, limit: 1 },
    });

    const { data: completedEvents } = useGetEvents({
        status: "done",
        pagination: { page: 1, limit: 1 },
    });

    const { data: rejectedUpdates } = useGetAllUpdatedEvents({
        status: "rejected",
        pagination: { page: 1, limit: 1 },
    });

    const { data: rejectedEvents } = useGetEvents({
        status: "rejected",
        pagination: { page: 1, limit: 1 },
    });

    const activeCount = activeEvents?.total_rows ?? 0;
    const pendingCount = (pendingEvents?.total_rows ?? 0) + (pendingUpdates?.total_rows ?? 0);
    const completedCount = completedEvents?.total_rows ?? 0;
    const rejectedCount = (rejectedEvents?.total_rows ?? 0) + (rejectedUpdates?.total_rows ?? 0);
    const totalEvents = activeCount + pendingCount + completedCount + rejectedCount;

    const pct = (val: number) =>
        totalEvents > 0 ? Math.round((val / totalEvents) * 100) + "%" : "0%";

    const navigate = useNavigate();

    const statCards = [
        {
            label: "Active events",
            value: activeCount,
            accent: "bg-teal-500",
            icon: <RiCalendarEventFill className="text-xl text-gray-400" />,
            route: "/admin/events/active",
        },
        {
            label: "Pending review",
            value: pendingCount,
            accent: "bg-amber-500",
            icon: <RiTimeFill className="text-xl text-gray-400" />,
            route: "/admin/events/pending",
        },
        {
            label: "Completed",
            value: completedCount,
            accent: "bg-gray-400",
            icon: <RiCalendarCheckFill className="text-xl text-gray-400" />,
            route: "/admin/events/done",
        },
        {
            label: "Rejected",
            value: rejectedCount,
            accent: "bg-red-500",
            icon: <MdCancel className="text-xl text-gray-400" />,
            route: "/admin/events/rejected",
        },
    ];

    const breakdownRows = [
        { label: "Active", dot: "bg-teal-500", val: pct(activeCount) },
        { label: "Pending", dot: "bg-amber-500", val: pct(pendingCount) },
        { label: "Completed", dot: "bg-gray-400", val: pct(completedCount) },
        { label: "Rejected", dot: "bg-red-500", val: pct(rejectedCount) },
    ];

    return (
        <AdminSidebar>
            <div className="p-6 w-full">

                {/* Header */}
                <div className="flex items-start justify-between flex-wrap gap-3 mb-8">
                    <div>
                        <h1 className="text-2xl font-semibold text-gray-900">Hello, Admin</h1>
                        <p className="text-sm text-gray-500 mt-1">
                            Here's what's happening across the platform.
                        </p>
                    </div>
                    {pendingCount > 0 && (
                        <span className="inline-flex items-center gap-2 text-xs bg-amber-50 text-amber-700 border border-amber-200 px-3 py-2 rounded-lg">
                            <span className="w-2 h-2 rounded-full bg-amber-500 animate-pulse" />
                            {pendingCount} item{pendingCount !== 1 ? "s" : ""} need review
                        </span>
                    )}
                </div>

                {/* Events stat cards */}
                <p className="text-xs font-medium tracking-widest text-gray-400 uppercase mb-3">
                    Events
                </p>
                <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 mb-6">
                    {statCards.map((card) => (
                        <button
                            key={card.label}
                            onClick={() => navigate(card.route)}
                            className="relative bg-white border border-gray-200 rounded-xl p-4 text-left hover:border-gray-300 hover:shadow-sm transition-all duration-150 overflow-hidden group"
                        >
                            <div
                                className={`absolute left-0 top-0 bottom-0 w-1 rounded-l-xl ${card.accent}`}
                            />
                            <p className="text-xs text-gray-500 mb-3 pl-1">{card.label}</p>
                            <div className="flex items-end justify-between pl-1">
                                <span className="text-3xl font-semibold text-gray-900 leading-none">
                                    {card.value}
                                </span>
                                <span className="group-hover:text-gray-600 transition-colors">
                                    {card.icon}
                                </span>
                            </div>
                        </button>
                    ))}
                </div>

                {/* Lower panels */}
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mb-4">
                    {/* Users panel */}
                    <div className="bg-white border border-gray-200 rounded-xl p-5">
                        <p className="text-sm font-medium text-gray-800 mb-4 flex items-center gap-2">
                            <FaRegUser className="text-gray-400" />
                            Users
                        </p>

                        <button
                            onClick={() => navigate("/admin/admin-list")}
                            className="w-full flex items-center gap-3 py-3 border-b border-gray-100 hover:bg-gray-50 rounded-lg px-2 -mx-2 transition-colors"
                        >
                            <div className="w-9 h-9 rounded-full bg-blue-50 flex items-center justify-center text-xs font-medium text-blue-700 flex-shrink-0">
                                AD
                            </div>
                            <div className="flex-1 text-left">
                                <p className="text-sm font-medium text-gray-800">Admins</p>
                                <p className="text-xs text-gray-400">Full platform access</p>
                            </div>
                            <span className="text-xl font-semibold text-gray-900">
                                {admins.isLoading ? (
                                    <span className="text-sm text-gray-400">...</span>
                                ) : (
                                    admins.data?.total_rows ?? 0
                                )}
                            </span>
                        </button>

                        <button
                            onClick={() => navigate("/admin/event-owners")}
                            className="w-full flex items-center gap-3 py-3 hover:bg-gray-50 rounded-lg px-2 -mx-2 transition-colors mt-1"
                        >
                            <div className="w-9 h-9 rounded-full bg-purple-50 flex items-center justify-center text-xs font-medium text-purple-700 flex-shrink-0">
                                EO
                            </div>
                            <div className="flex-1 text-left">
                                <p className="text-sm font-medium text-gray-800">Event owners</p>
                                <p className="text-xs text-gray-400">Can submit & manage events</p>
                            </div>
                            <span className="text-xl font-semibold text-gray-900">
                                {organizers.isLoading ? (
                                    <span className="text-sm text-gray-400">...</span>
                                ) : (
                                    organizers.data?.total_rows ?? 0
                                )}
                            </span>
                        </button>
                    </div>

                    {/* Breakdown panel */}
                    <div className="bg-white border border-gray-200 rounded-xl p-5">
                        <p className="text-sm font-medium text-gray-800 mb-4">Event breakdown</p>
                        <div className="space-y-0">
                            {breakdownRows.map((row, i) => (
                                <div
                                    key={row.label}
                                    className={`flex items-center justify-between py-2.5 ${i < breakdownRows.length - 1 ? "border-b border-gray-100" : ""
                                        }`}
                                >
                                    <span className="flex items-center gap-2 text-sm text-gray-500">
                                        <span className={`w-2 h-2 rounded-full ${row.dot}`} />
                                        {row.label}
                                    </span>
                                    <span className="text-sm font-medium text-gray-800">
                                        {row.val} of total
                                    </span>
                                </div>
                            ))}
                        </div>
                    </div>
                </div>

                {/* CTA bar */}
                {pendingCount > 0 && (
                    <div className="bg-gray-50 border border-gray-200 rounded-xl px-5 py-4 flex items-center justify-between gap-4 flex-wrap">
                        <p className="text-sm text-gray-500">
                            <span className="font-medium text-gray-900">
                                {pendingCount} event{pendingCount !== 1 ? "s" : ""}
                            </span>{" "}
                            are waiting for your review — approvals keep owners moving.
                        </p>
                        <button
                            onClick={() => navigate("/admin/events/pending")}
                            className="inline-flex items-center gap-2 text-sm px-4 py-2 rounded-lg border border-gray-300 bg-white hover:bg-gray-50 text-gray-700 transition-colors font-medium flex-shrink-0"
                        >
                            Review pending
                            <span>→</span>
                        </button>
                    </div>
                )}
            </div>
        </AdminSidebar>
    );
}