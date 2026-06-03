import type { Dispatch, SetStateAction } from "react";
import type { DateFilterType } from "../../utils/dateFilter";
import FilterButton from "./filterButton";
import DateSelector from "./dateSelector";
import MonthSelector from "./monthSelector";
import YearSelector from "./yearSelector";


interface Props {
    dateFilterType: DateFilterType;
    setDateFilterType:
    Dispatch<SetStateAction<DateFilterType>>;

    selectedDate: Date | null;
    setSelectedDate:
    Dispatch<SetStateAction<Date | null>>;

    selectedMonth?: number;
    setSelectedMonth:
    Dispatch<SetStateAction<number | undefined>>;

    selectedYear?: number;
    setSelectedYear:
    Dispatch<SetStateAction<number | undefined>>;
}

export default function DateFilterDropdown({
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
        <div className="flex flex-col gap-3">
            <FilterButton
                value="all"
                label="Any Time"
                activeValue={dateFilterType}
                onClick={setDateFilterType}
            />

            <FilterButton
                value="today"
                label="Today"
                activeValue={dateFilterType}
                onClick={setDateFilterType}
            />

            <FilterButton
                value="this_month"
                label="This Month"
                activeValue={dateFilterType}
                onClick={setDateFilterType}
            />

            <FilterButton
                value="this_year"
                label="This Year"
                activeValue={dateFilterType}
                onClick={setDateFilterType}
            />

            <hr />

            <FilterButton
                value="date"
                label="Specific Date"
                activeValue={dateFilterType}
                onClick={setDateFilterType}
            />

            {dateFilterType === "date" && (
                <DateSelector
                    value={selectedDate}
                    onChange={setSelectedDate}
                />
            )}

            <FilterButton
                value="month"
                label="Month"
                activeValue={dateFilterType}
                onClick={setDateFilterType}
            />

            {dateFilterType === "month" && (
                <MonthSelector
                    month={selectedMonth}
                    year={selectedYear}
                    onMonthChange={setSelectedMonth}
                    onYearChange={setSelectedYear}
                />
            )}

            <FilterButton
                value="year"
                label="Year"
                activeValue={dateFilterType}
                onClick={setDateFilterType}
            />

            {dateFilterType === "year" && (
                <YearSelector
                    year={selectedYear}
                    onChange={setSelectedYear}
                />
            )}
        </div>
    );
}