interface Props {
    year?: number;
    onChange: (year: number | undefined) => void;
}

export default function YearSelector({
    year,
    onChange,
}: Props) {
    const years = Array.from(
        { length: 10 },
        (_, i) => new Date().getFullYear() - 5 + i
    );

    return (
        <select
            value={year ?? ""}
            onChange={(e) =>
                onChange(
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
    );
}