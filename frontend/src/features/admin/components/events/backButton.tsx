import { MdKeyboardArrowLeft } from "react-icons/md";

interface BackButtonProps {
    onClick?: () => void;
    label?: string;
}

export default function BackButton({ onClick, label = "back" }: BackButtonProps) {
    return (
        <button
            onClick={onClick}
            className="inline-flex items-center gap-1 text-sm text-gray-500 hover:text-gray-800 transition-colors mb-5 group"
        >
            <MdKeyboardArrowLeft
                className="w-4 h-4 transition-transform group-hover:-translate-x-0.5"
                size={25}
            />
            {label}
        </button>
    )
}