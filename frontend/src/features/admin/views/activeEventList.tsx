import PageHeader from "../components/events/pageHeader";
import AdminSidebar from "../components/sideBar";
import EventList from "../components/events/eventList";
import { useState } from "react";

export default function ActiveEventList() {
    const [currentPage, setCurrentPage] = useState(1);
    const totalPage = 4;
    return (
        <AdminSidebar>
            <>
                <PageHeader
                    title="Active Events"
                    description="Monitor and manage all currently live events across the platform."
                    pendingPage={false}
                />

                <EventList
                    pendingEvent={false}
                    currentPage={currentPage}
                    setCurrentPage={setCurrentPage}
                    totalPage={totalPage}
                />
            </>
        </AdminSidebar>
    )
}