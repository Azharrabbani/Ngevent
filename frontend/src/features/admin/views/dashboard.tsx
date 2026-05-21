import { FaRegUser } from "react-icons/fa";
import { RiCalendarEventFill } from "react-icons/ri";
import { useListUsers } from "../../profile/hooks/useListUsers";
import AdminSidebar from "../components/sideBar";
import { useNavigate } from "react-router-dom";
import { useListAttendee } from "../../profile/hooks/attendee/useListAttendee";
import { useGetEvents } from "../../events/hooks/useGetEvents";

export default function AdminDashboard() {
    const { data: attendees, isLoading: attendessLoading } = useListAttendee({
        filter: "",
    })

    const organizers = useListUsers({
        role: "event organizer"
    })

    const { data: activeEvents } = useGetEvents({
        status: "active",
    })

    const { data: pendingEvents } = useGetEvents({
        status: "pending",
    })

    const { data: completedEvents } = useGetEvents({
        status: "done",
    })

    const { data: rejectedEvent } = useGetEvents({
        status: "rejected",
    })

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
                            <h1 className="text-2xl font-bold">{pendingEvents?.total_rows ?? 0}</h1>
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
                            <h1 className="text-2xl font-bold">{rejectedEvent?.total_rows ?? 0}</h1>
                        </div>
                    </div>

                    <div
                        className="bg-blue-500 rounded-xl p-4 text-white hover:-translate-y-1 hover:scale-110 transition duration-300 cursor-pointer ease-in-out"
                        onClick={() => navigate("/admin/attendee-list")}
                    >
                        <h1 className="text-lg">Attendees</h1>
                        <div className="flex items-center justify-between mt-4">
                            <FaRegUser className="text-4xl" />
                            {attendessLoading && <h1>Loading...</h1>}
                            <h1 className="text-2xl font-bold">{attendees?.total_rows}</h1>
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