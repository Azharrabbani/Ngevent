import { useState, useEffect, useRef } from "react";
import { useGetEventsActive } from "../../hooks/useGeEventsActive"
import DashboardHeader from "../components/header/dashboardHeader";
import EventCreatorTabs from "../components/tabs/eventCreatorTabs";
import CategorySlider from "../components/category/categorySlider";
import { buildEventDateFilters, type DateFilterType } from "../../../../utils/dateFilter";
import FilterSection from "../components/dashboard/filterSection";
import EventGrid from "../components/event/eventGrid";
import { useUserLocation } from "../../hooks/useUserLocation";
import { useSearchParams } from "react-router-dom";
import { FiMap, FiList } from "react-icons/fi";
import EventsMapContainer from "../components/map/eventMapContainer";

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
    const [searchParams, setSearchParams] = useSearchParams();
    const [viewMode, setViewMode] = useState<"list" | "map">("list");

    const [selectedCategory, setSelectedCategory] = useState<number[] | undefined>(
        searchParams.get("category")
            ? [Number(searchParams.get("category"))]
            : undefined
    );

    const [location, setLocation] = useState(
        searchParams.get("location") || undefined
    );

    const [dateFilterType, setDateFilterType] =
        useState<DateFilterType>(
            (searchParams.get("date") as DateFilterType) || "all"
        );

    const [selectedDate, setSelectedDate] =
        useState<Date | null>(
            searchParams.get("selectedDate")
                ? new Date(
                    searchParams.get(
                        "selectedDate"
                    )!
                )
                : null
        );

    const [selectedMonth, setSelectedMonth] =
        useState<number | undefined>(
            searchParams.get("month")
                ? Number(
                    searchParams.get("month")
                )
                : undefined
        );

    const [selectedYear, setSelectedYear] =
        useState<number | undefined>(
            searchParams.get("year")
                ? Number(
                    searchParams.get("year")
                )
                : undefined
        );

    const [search, setSearch] = useState(
        searchParams.get("search") || ""
    );

    const [nearestEnabled, setNearestEnabled] = useState(
        searchParams.get("nearest") === "true"
    );

    const { lat, lon, loading: locationLoading, denied } = useUserLocation();

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
    } = useGetEventsActive(
        {
            category: selectedCategory,
            location: location,
            search: search || undefined,
            ...dateFilters,
            limit: 8,
            ...(nearestEnabled && lat && lon ? { lat, lon } : {}),
        },
        viewMode === "list"
    );

    const events = data?.pages.flatMap((page) => page.rows ?? []) || [];
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

    useEffect(() => {
        const params = new URLSearchParams();

        if (search)
            params.set("search", search);

        if (location)
            params.set("location", location);

        if (selectedCategory?.length)
            params.set(
                "category",
                String(selectedCategory[0])
            );

        if (selectedDate)
            params.set(
                "selectedDate",
                selectedDate.toISOString()
            );

        if (selectedMonth)
            params.set(
                "month",
                String(selectedMonth)
            );

        if (selectedYear)
            params.set(
                "year",
                String(selectedYear)
            );

        if (nearestEnabled)
            params.set("nearest", "true");

        if (dateFilterType !== "all")
            params.set("date", dateFilterType);

        setSearchParams(params, {
            replace: true,
        });
    }, [
        search,
        location,
        selectedCategory,
        nearestEnabled,
        dateFilterType,
        setSearchParams,
    ]);

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

                    <div className="flex items-center gap-3">
                        <button
                            onClick={() => setViewMode((prev) => (prev === "list" ? "map" : "list"))}
                            className="flex items-center gap-2 px-4 py-2.5 rounded-xl border border-slate-200 bg-white text-sm font-medium text-slate-700 hover:border-slate-300 transition"
                        >
                            {viewMode === "list" ? (
                                <>
                                    <FiMap className="w-4 h-4" />
                                    View On The Map
                                </>
                            ) : (
                                <>
                                    <FiList className="w-4 h-4" />
                                    View On The List
                                </>
                            )}
                        </button>
                        <EventCreatorTabs activeTab="event" />
                    </div>
                </div>

                <CategorySlider
                    selectedCategory={selectedCategory}
                    onChange={setSelectedCategory}
                />

                <div className="mt-6">
                    <FilterSection
                        viewMode={viewMode}
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
                        nearestEnabled={nearestEnabled}
                        setNearestEnabled={setNearestEnabled}
                        locationLoading={locationLoading}
                        locationDenied={denied}
                    />
                </div>

                {viewMode === "map" ? (
                    <div className="mt-8">
                        <EventsMapContainer
                            category={selectedCategory}
                            location={location}
                            search={search || undefined}
                            dateFilters={dateFilters}
                        />
                    </div>
                ) : isLoading ? (
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