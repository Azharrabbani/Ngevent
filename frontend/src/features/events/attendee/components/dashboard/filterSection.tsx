import type { DateFilterType } from "../../../../../utils/dateFilter";
import EventDateFilter from "./eventDateFilter";
import LocationFilter from "./locationFilter";


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
        </div>
    );
}