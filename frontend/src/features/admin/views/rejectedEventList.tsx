import { useEffect, useState } from "react";
import { useGetEvents } from "../../events/hooks/useGetEvents";
import AdminSidebar from "../components/sideBar";
import EventList from "../components/events/eventList";

export default function RejectedEventList() {
    const [currentPage, setCurrentPage] = useState(1);
    const [search, setSearch] = useState<string | undefined>(undefined);
    const [sort, setSort] = useState<string | undefined>("desc");
    const [dateFilter, setDateFilter] = useState<string | undefined>(undefined);
    const [getUpdate, setGetUpdate] = useState<boolean | undefined>(undefined);

    const { data, isLoading } = useGetEvents({
        status: "rejected",
        search: search,
        sort: sort,
        date: dateFilter,
        get_update: getUpdate,
        pagination: { limit: 4, page: currentPage },
    });

    const totalPage = data?.total_pages ?? 1;

    useEffect(() => {
        const delay = setTimeout(() => setCurrentPage(1), 500);
        return () => clearTimeout(delay);
    }, [search]);
    return (
        <AdminSidebar>
            <>
                <EventList
                    status="rejected"
                    data={data}
                    isLoading={isLoading}
                    reviewEvent={true}
                    currentPage={currentPage}
                    totalPage={totalPage}
                    search={search}
                    setSearch={setSearch}
                    sort={sort}
                    setSort={setSort}
                    dateFilter={dateFilter}
                    setDateFilter={setDateFilter}
                    getUpdate={getUpdate}
                    setGetupdate={setGetUpdate}
                    setCurrentPage={setCurrentPage}
                />
            </>
        </AdminSidebar>
    )
}