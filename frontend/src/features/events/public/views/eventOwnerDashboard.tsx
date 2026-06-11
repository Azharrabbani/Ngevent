import { useState, useEffect, useRef } from "react";
import { useNavigate } from "react-router-dom";

import EventCreatorTabs from "../components/tabs/eventCreatorTabs";
import OrganizerSkeletonCard from "../components/skeleton/eventOwnerSkeleton";
import { useGetPublicOrganizers } from "../../../profile/hooks/organizer/useGetPublicOrganizers";
import type { OrganizerResponse } from "../../../profile/types/profileResponse";
import DashboardHeader from "../components/header/dashboardHeader";
import EventOwnerGrid from "../components/owner/eventOwnerGrid";

export default function PublicEventOwnerDashboard() {
    const [search, setSearch] = useState<string>("");
    const navigate = useNavigate();

    const {
        data,
        isLoading,
        isFetchingNextPage,
        hasNextPage,
        fetchNextPage,
    } = useGetPublicOrganizers({
        filter: search || undefined,
        limit: 9,
    });

    const organizers: OrganizerResponse[] =
        data?.pages.flatMap((page) => page.rows) || [];

    const observerTarget = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const observer = new IntersectionObserver(
            (entries) => {
                if (
                    entries[0].isIntersecting &&
                    hasNextPage &&
                    !isFetchingNextPage
                ) {
                    fetchNextPage();
                }
            },
            { threshold: 0.1 }
        );

        const currentTarget = observerTarget.current;
        if (currentTarget) {
            observer.observe(currentTarget);
        }

        return () => {
            if (currentTarget) {
                observer.unobserve(currentTarget);
            }
        };
    }, [hasNextPage, isFetchingNextPage, fetchNextPage]);

    const handleOrganizerSelect = (organizer: OrganizerResponse) => {
        navigate(`/event-creators/${organizer.id}`);
    };

    return (
        <div className="min-h-screen bg-slate-50 flex flex-col">
            <DashboardHeader onSearchChange={setSearch} />

            <main className="flex-1 max-w-7xl mx-auto w-full px-4 lg:px-6 py-10">
                <div className="flex items-start justify-between flex-wrap gap-4">
                    <div>
                        <h1 className="text-3xl font-bold text-slate-900">
                            Find your{" "}
                            <span className="text-transparent bg-clip-text bg-gradient-to-r from-blue-600 to-indigo-500">
                                Favorite Event Owner
                            </span>
                        </h1>
                        <p className="mt-1 text-slate-500 text-sm">
                            Discover the best event owner around you.
                        </p>
                    </div>

                    <EventCreatorTabs activeTab="event_owner" />
                </div>

                {isLoading ? (
                    <div className="mt-8 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                        {Array.from({ length: 9 }).map((_, i) => (
                            <OrganizerSkeletonCard key={i} />
                        ))}
                    </div>
                ) : (
                    <>
                        <EventOwnerGrid
                            organizers={organizers}
                            onSelect={handleOrganizerSelect}
                        />

                        <div
                            ref={observerTarget}
                            className="h-20 flex items-center justify-center mt-4"
                        >
                            {isFetchingNextPage && (
                                <div className="flex items-center gap-2 text-blue-600 font-medium">
                                    <div className="w-5 h-5 border-2 border-blue-600 border-t-transparent rounded-full animate-spin" />
                                    <span>Loading more owners...</span>
                                </div>
                            )}
                        </div>
                    </>
                )}
            </main>
        </div>
    );
}