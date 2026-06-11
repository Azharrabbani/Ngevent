import { useState, useEffect, useRef } from "react";
import { useGetEventsActive } from "../../hooks/useGeEventsActive"
import DashboardHeader from "../components/header/dashboardHeader";
import EventCreatorTabs from "../components/tabs/eventCreatorTabs";
import CategorySlider from "../components/category/categorySlider";
import { buildEventDateFilters, type DateFilterType } from "../../../../utils/dateFilter";
import FilterSection from "../components/dashboard/filterSection";
import EventGrid from "../components/event/eventGrid";

function SkeletonCard() {
    return (
        <div className="bg-white rounded-2xl overflow-hidden border border-slate-200 animate-pulse">
            <div className="h-48 bg-slate-200" />
            <div className="p-4 space-y-3">
                <div className="h-4 bg-slate-200 rounded w-3/4" />
                <div className="h-3 bg-slate-100 rounded w-1/2" />
                <div className="h-3 bg-slate-100 rounded w-2/5" />
                <div className="pt-3 border-t border-slate-100">
                    <div className="h-3 bg-slate-100 rounded w-1/3" />
                </div>
            </div>
        </div>
    );
}

export default function PublicDashboard() {
    const [selectedCategory, setSelectedCategory] = useState<number[]>();
    const [location, setLocation] = useState<string>();
    const [dateFilterType, setDateFilterType] = useState<DateFilterType>("all");
    const [selectedDate, setSelectedDate] = useState<Date | null>(null);
    const [selectedMonth, setSelectedMonth] = useState<number>();
    const [selectedYear, setSelectedYear] = useState<number>();
    const [search, setSearch] = useState<string>("");

    const dateFilters = buildEventDateFilters(
        dateFilterType,
        selectedDate,
        selectedMonth,
        selectedYear
    );

    const {
        data,
        isLoading,
        isFetchingNextPage,
        hasNextPage,
        fetchNextPage
    } = useGetEventsActive({
        category: selectedCategory,
        location: location,
        search: search || undefined,
        ...dateFilters,
        limit: 8
    });

    const events = data?.pages.flatMap((page) => page.rows) || [];
    const observerTarget = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const observer = new IntersectionObserver(
            (entries) => {
                if (entries[0].isIntersecting && hasNextPage && !isFetchingNextPage) {
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

    return (
        <div className="min-h-screen bg-slate-50">
            <DashboardHeader onSearchChange={setSearch} />

            <main className="max-w-7xl mx-auto px-4 lg:px-6 py-10">
                <div className="flex items-center justify-between flex-wrap gap-4">
                    <div>
                        <h1 className="text-3xl font-bold text-slate-900">
                            Find your{" "}
                            <span className="text-transparent bg-clip-text bg-gradient-to-r from-blue-600 to-indigo-500">
                                favorite event
                            </span>
                        </h1>
                        <p className="mt-1 text-slate-500 text-sm">
                            Discover the most exciting events around you
                        </p>
                    </div>
                    <EventCreatorTabs activeTab="event" />
                </div>

                <CategorySlider
                    selectedCategory={selectedCategory}
                    onChange={setSelectedCategory}
                />

                <div className="mt-6">
                    <FilterSection
                        location={location}
                        setLocation={setLocation}
                        dateFilterType={dateFilterType}
                        setDateFilterType={setDateFilterType}
                        selectedDate={selectedDate}
                        setSelectedDate={setSelectedDate}
                        selectedMonth={selectedMonth}
                        setSelectedMonth={setSelectedMonth}
                        selectedYear={selectedYear}
                        setSelectedYear={setSelectedYear}
                    />
                </div>

                {isLoading ? (
                    <div className="mt-8 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-5">
                        {Array.from({ length: 8 }).map((_, i) => (
                            <SkeletonCard key={i} />
                        ))}
                    </div>
                ) : (
                    <>
                        <EventGrid events={events} />

                        <div ref={observerTarget} className="h-20 flex items-center justify-center mt-8">
                            {isFetchingNextPage && (
                                <div className="flex items-center gap-2 text-blue-600 font-medium">
                                    <div className="w-5 h-5 border-2 border-blue-600 border-t-transparent rounded-full animate-spin" />
                                    <span>Loading more events...</span>
                                </div>
                            )}
                        </div>
                    </>
                )}
            </main>
        </div>
    );
}
