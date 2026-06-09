interface Props {
    label: string;
    active?: boolean;
    onClick: () => void;
}

export default function CategoryItem({
    label,
    active,
    onClick,
}: Props) {
    return (
        <button
            onClick={onClick}
            className={`
                shrink-0
                px-5 py-2
                rounded-full
                text-sm font-medium
                border
                transition-all duration-200
                whitespace-nowrap
                ${active
                    ? "bg-blue-600 text-white border-blue-600 shadow-md shadow-blue-200"
                    : "bg-white text-slate-600 border-slate-200 hover:border-blue-400 hover:text-blue-600"
                }
            `}
        >
            {label}
        </button>
    );
}