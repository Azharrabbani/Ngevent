import { useCallback, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import DashboardHeader from "../components/header/dashboardHeader";
import { useGetOrganizerBySlug } from "../../../profile/hooks/organizer/useGetOrganizerBySlug";
import { useGetPublicOrganizerEvents } from "../../hooks/useGetPublicOrganizerEvents";
import OwnerProfileSkeleton from "../components/skeleton/ownerProfileSkeleton";
import OwnerProfileCard from "../components/event_owner/ownerProfileCard";
import EventStatusTabs from "../components/tabs/eventStatusTabs";
import EventCardSkeleton from "../components/skeleton/eventCardSkeleton";
import PaginationTabs from "../components/tabs/paginationTabs";
import { ShockIcon } from "../../../../components/icon";
import { defaultPagination } from "../../../../utils/pagination";
import OwnerEventCard from "../components/event/ownerEventCard";

type EventStatus = "active" | "done";

const ITEMS_PER_PAGE = 8;

function buildShowingLabel(page: number, limit: number, total: number): string {
    if (total === 0) return "No events found";
    const from = (page - 1) * limit + 1;
    const to = Math.min(page * limit, total);
    return `Showing ${from}–${to} of a total of ${total} events`;
}

export default function OwnerViewPage() {
    const { slug } = useParams<{ slug: string }>();
    const navigate = useNavigate();

    const [search, setSearch] = useState<string>("");
    const [eventStatus, setEventStatus] = useState<EventStatus>("active");
    const [currentPage, setCurrentPage] = useState<number>(1);

    const { data: owner, isLoading: ownerLoading } = useGetOrganizerBySlug(slug!);

    const { data: eventsData, isLoading: eventsLoading } = useGetPublicOrganizerEvents(owner?.id ?? "", {
        status: eventStatus === "active" ? "active" : "done",
        title: search || undefined,
        pagination: defaultPagination(currentPage)
    });

    const events = eventsData?.rows ?? [];
    const totalEvents = eventsData?.total_rows ?? 0;
    const totalPages = eventsData?.total_pages ?? 1;

    const handleStatusChange = (status: EventStatus) => {
        setEventStatus(status);
        setCurrentPage(1);
    };

    const handleSearchChange = useCallback((val: string) => {
        setSearch(val);
        setCurrentPage(1);
    }, []);

    return (
        <div className="min-h-screen bg-slate-50 flex flex-col">
            <DashboardHeader
                onSearchChange={handleSearchChange}
                isEventOwnerDashboard={true}
            />

            <main className="flex-1 max-w-7xl mx-auto w-full px-4 lg:px-6 py-10 space-y-8">
                {ownerLoading || !owner ? (
                    <OwnerProfileSkeleton />
                ) : (
                    <OwnerProfileCard organizer={owner} />
                )}

                <section>
                    <div className="flex flex-wrap items-center justify-between gap-4">
                        <EventStatusTabs
                            activeStatus={eventStatus}
                            onChange={handleStatusChange}
                        />

                        {!eventsLoading && (
                            <p className="text-sm text-slate-500">
                                {buildShowingLabel(currentPage, ITEMS_PER_PAGE, totalEvents)}
                            </p>
                        )}
                    </div>

                    {eventsLoading ? (
                        <div className="mt-6 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5">
                            {Array.from({ length: ITEMS_PER_PAGE }).map((_, i) => (
                                <EventCardSkeleton key={i} />
                            ))}
                        </div>
                    ) : events.length === 0 ? (
                        <div className="mt-16 flex flex-col items-center text-center">
                            <div className="w-14 h-14 rounded-2xl bg-slate-100 flex items-center justify-center mb-3">
                                <ShockIcon className="text-5xl text-indigo-500" />
                            </div>
                            <p className="text-slate-700 font-semibold">No events found</p>
                            <p className="text-slate-400 text-sm mt-1">
                                {eventStatus === "active"
                                    ? "This owner has no active events right now."
                                    : "No past events to display."}
                            </p>
                        </div>
                    ) : (
                        <div className="mt-6 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5">
                            {events.map((event) => (
                                <OwnerEventCard
                                    key={event.id}
                                    event={event}
                                    onClick={() => navigate(`/events/${event.event.slug}`)}
                                />
                            ))}
                        </div>
                    )}

                    {!eventsLoading && totalPages > 1 && (
                        <PaginationTabs
                            currentPage={currentPage}
                            totalPages={totalPages}
                            onPageChange={(page) => {
                                setCurrentPage(page);
                                window.scrollTo({ top: 0, behavior: "smooth" });
                            }}
                        />
                    )}
                </section>
            </main>
        </div>
    );
}