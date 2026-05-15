import { useState } from "react";
import EventList from "../components/events/eventList";
import PageHeader from "../components/events/pageHeader";
import AdminSidebar from "../components/sideBar";

export default function PendingEventList() {
    const [currentPage, setCurrentPage] = useState(1);
    const totalPage = 4;
    return (
        <AdminSidebar>
            <>
                <PageHeader
                    title="Pending Events"
                    description="Review and manage event submissions requiring approval."
                    pendingPage={true}
                />

                <EventList
                    pendingEvent={true}
                    currentPage={currentPage}
                    totalPage={totalPage}
                    setCurrentPage={setCurrentPage}
                />
            </>
        </AdminSidebar>
    )
}