import type { DateFilterType } from "../../../../../utils/dateFilter";
import EventDateFilter from "./eventDateFilter";
import LocationFilter from "./locationFilter";
import NearestFilter from "./nearestFilter";

interface Props {
    location?: string;
    setLocation: (value?: string) => void;
    dateFilterType: DateFilterType;
    setDateFilterType: React.Dispatch<React.SetStateAction<DateFilterType>>;
    selectedDate: Date | null;
    setSelectedDate: React.Dispatch<React.SetStateAction<Date | null>>;
    selectedMonth?: number;
    setSelectedMonth: React.Dispatch<React.SetStateAction<number | undefined>>;
    selectedYear?: number;
    setSelectedYear: React.Dispatch<React.SetStateAction<number | undefined>>;
    // Nearest filter
    nearestEnabled: boolean;
    setNearestEnabled: React.Dispatch<React.SetStateAction<boolean>>;
    locationLoading: boolean;
    locationDenied: boolean;
}

export default function FilterSection({
    location,
    setLocation,
    dateFilterType,
    setDateFilterType,
    selectedDate,
    setSelectedDate,
    selectedMonth,
    setSelectedMonth,
    selectedYear,
    setSelectedYear,
    nearestEnabled,
    setNearestEnabled,
    locationLoading,
    locationDenied,
}: Props) {
    return (
        <div className="flex flex-wrap items-center gap-3">
            <LocationFilter
                value={location}
                onChange={setLocation}
            />

            <EventDateFilter
                dateFilterType={dateFilterType}
                setDateFilterType={setDateFilterType}
                selectedDate={selectedDate}
                setSelectedDate={setSelectedDate}
                selectedMonth={selectedMonth}
                setSelectedMonth={setSelectedMonth}
                selectedYear={selectedYear}
                setSelectedYear={setSelectedYear}
            />

            <NearestFilter
                enabled={nearestEnabled}
                onToggle={setNearestEnabled}
                locationLoading={locationLoading}
                denied={locationDenied}
            />
        </div>
    );
}