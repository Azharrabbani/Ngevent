import { FaRegUser } from "react-icons/fa";
import { RiCalendarEventFill } from "react-icons/ri";
import { useListUsers } from "../../profile/hooks/useListUsers";
import AdminSidebar from "../components/sideBar";
import { useNavigate } from "react-router-dom";
import { useGetEvents } from "../../events/hooks/useGetEvents";
import { useGetAllUpdatedEvents } from "../../events/hooks/useGetAllUpdatedEvent";

export default function AdminDashboard() {
    const organizers = useListUsers({
        role: "event organizer",
        pagination: { page: 1, limit: 1 }
    })

    const admins = useListUsers({
        role: "admin",
        pagination: { page: 1, limit: 1 }
    })

    const { data: activeEvents } = useGetEvents({
        status: "active",
        pagination: { page: 1, limit: 1 }
    })

    const { data: pendingUpdates } = useGetAllUpdatedEvents({
        status: "pending",
        pagination: { page: 1, limit: 1 },
    });

    const { data: pendingEvents } = useGetEvents({
        status: "pending",
        pagination: { page: 1, limit: 1 },
    });

    const pendingEventCount = pendingEvents?.total_rows ?? 0;
    const pendingUpdateCount = pendingUpdates?.total_rows ?? 0;
    const totalPending = pendingEventCount + pendingUpdateCount


    const { data: completedEvents } = useGetEvents({
        status: "done",
    })

    const { data: rejectedUpdates } = useGetAllUpdatedEvents({
        status: "rejected",
        pagination: { page: 1, limit: 1 },
    });

    const { data: rejectedEvents } = useGetEvents({
        status: "rejected",
        pagination: { page: 1, limit: 1 },
    });

    const rejectedEventCount = rejectedEvents?.total_rows ?? 0;
    const rejectedUpdateCount = rejectedUpdates?.total_rows ?? 0;
    const totalRejected = rejectedEventCount + rejectedUpdateCount;

    const navigate = useNavigate();


    return (
        <AdminSidebar>
            <>
                <h1 className="text-2xl font-semibold">Hello Admin!</h1>

                <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-6 p-10 gap-6">
                    <div
                        className="bg-green-500 rounded-xl p-4 text-white hover:-translate-y-1 hover:scale-110 transition duration-300 cursor-pointer ease-in-out"
                        onClick={() => navigate("/admin/events/active")}
                    >
                        <h1 className="text-lg">Active Events</h1>
                        <div className="flex items-center justify-between mt-4">
                            <RiCalendarEventFill className="text-4xl" />
                            <h1 className="text-2xl font-bold">{activeEvents?.total_rows ?? 0}</h1>
                        </div>
                    </div>

                    <div
                        className="bg-amber-500 rounded-xl p-4 text-white hover:-translate-y-1 hover:scale-110 transition duration-300 cursor-pointer ease-in-out"
                        onClick={() => navigate("/admin/events/pending")}
                    >
                        <h1 className="text-lg">Pending Events</h1>
                        <div className="flex items-center justify-between mt-4">
                            <RiCalendarEventFill className="text-4xl" />
                            <h1 className="text-2xl font-bold">{totalPending ?? 0}</h1>
                        </div>
                    </div>

                    <div
                        className="bg-gray-200 rounded-xl p-4 text-gray-600 hover:-translate-y-1 hover:scale-110 transition duration-300 cursor-pointer ease-in-out"
                        onClick={() => navigate("/admin/events/done")}
                    >
                        <h1 className="text-lg">Completed Events</h1>
                        <div className="flex items-center justify-between mt-4">
                            <RiCalendarEventFill className="text-4xl" />
                            <h1 className="text-2xl font-bold">{completedEvents?.total_rows ?? 0}</h1>
                        </div>
                    </div>

                    <div
                        className="bg-red-500 rounded-xl p-4 text-white hover:-translate-y-1 hover:scale-110 transition duration-300 cursor-pointer ease-in-out"
                        onClick={() => navigate("/admin/events/rejected")}
                    >
                        <h1 className="text-lg">Rejected Events</h1>
                        <div className="flex items-center justify-between mt-4">
                            <RiCalendarEventFill className="text-4xl" />
                            <h1 className="text-2xl font-bold">{totalRejected ?? 0}</h1>
                        </div>
                    </div>

                    <div
                        className="bg-blue-500 rounded-xl p-4 text-white hover:-translate-y-1 hover:scale-110 transition duration-300 cursor-pointer ease-in-out"
                        onClick={() => navigate("/admin/admin-list")}
                    >
                        <h1 className="text-lg">Admins</h1>
                        <div className="flex items-center justify-between mt-4">
                            <FaRegUser className="text-4xl" />
                            {admins.isLoading ? <h1>Loading...</h1> : null}
                            <h1 className="text-2xl font-bold">{admins.data?.total_rows}</h1>
                        </div>
                    </div>

                    <div
                        className="bg-[#312E81] rounded-xl p-4 text-white hover:-translate-y-1 hover:scale-110 transition duration-300 cursor-pointer ease-in-out"
                        onClick={() => navigate("/admin/organizer-list")}
                    >
                        <h1 className="text-lg">Event Organizers</h1>
                        <div className="flex items-center justify-between mt-4">
                            <FaRegUser className="text-4xl" />
                            <h1 className="text-2xl font-bold">{organizers.data?.total_rows}</h1>
                        </div>
                    </div>
                </div>
            </>
        </AdminSidebar>
    );
}