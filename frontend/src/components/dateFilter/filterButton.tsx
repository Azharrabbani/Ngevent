import type { DateFilterType } from "../../utils/dateFilter";


interface Props {
    value: DateFilterType;
    label: string;
    activeValue: DateFilterType;
    onClick: (value: DateFilterType) => void;
}

export default function FilterButton({
    value,
    label,
    activeValue,
    onClick,
}: Props) {
    return (
        <button
            onClick={() => onClick(value)}
            className={`text-left px-3 py-2 rounded transition
                ${activeValue === value
                    ? "bg-blue-100 text-blue-600"
                    : "hover:bg-gray-100"
                }`}
        >
            {label}
        </button>
    );
}