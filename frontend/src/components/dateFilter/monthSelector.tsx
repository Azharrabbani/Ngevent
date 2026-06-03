interface Props {
    month?: number;
    year?: number;
    onMonthChange: (month: number | undefined) => void;
    onYearChange: (year: number | undefined) => void;
}

export default function MonthSelector({
    month,
    year,
    onMonthChange,
    onYearChange,
}: Props) {
    const years = Array.from(
        { length: 10 },
        (_, i) => new Date().getFullYear() - 5 + i
    );

    const months = Array.from(
        { length: 12 },
        (_, i) => ({
            value: i + 1,
            label: new Date(0, i).toLocaleString(
                "en-US",
                {
                    month: "long",
                }
            ),
        })
    );

    return (
        <>
            <select
                value={month ?? ""}
                onChange={(e) =>
                    onMonthChange(
                        e.target.value
                            ? Number(e.target.value)
                            : undefined
                    )
                }
                className="border rounded px-3 py-2"
            >
                <option value="">
                    Select Month
                </option>

                {months.map((month) => (
                    <option
                        key={month.value}
                        value={month.value}
                    >
                        {month.label}
                    </option>
                ))}
            </select>

            <select
                value={year ?? ""}
                onChange={(e) =>
                    onYearChange(
                        e.target.value
                            ? Number(e.target.value)
                            : undefined
                    )
                }
                className="border rounded px-3 py-2"
            >
                <option value="">
                    Select Year
                </option>

                {years.map((year) => (
                    <option
                        key={year}
                        value={year}
                    >
                        {year}
                    </option>
                ))}
            </select>
        </>
    );
}