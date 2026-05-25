import { FiCheckCircle } from "react-icons/fi";
import { TiDeleteOutline } from "react-icons/ti";
import { IoCheckmark } from "react-icons/io5";
import { AiOutlineClose } from "react-icons/ai";

interface ApprovalButtonsProps {
    onApprove: () => Promise<void> | void;
    onReject: () => Promise<void> | void;
    variant?: "outline" | "solid";
    isSubmitting?: boolean;
}

export default function ApprovalButtons({ onApprove, onReject, variant = "solid", isSubmitting }: ApprovalButtonsProps) {
    if (variant === "outline") {
        return (
            <div className="flex items-center gap-3">
                <button
                    onClick={onApprove}
                    disabled={isSubmitting}
                    className="flex items-center gap-2 px-4 py-2 rounded-lg border border-gray-300 text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                    {isSubmitting ? (
                        <Spinner />
                    ) : (
                        <FiCheckCircle className="w-4 h-4 text-gray-500" size={20} />
                    )}
                    Approve
                </button>
                <button
                    onClick={onReject}
                    disabled={isSubmitting}
                    className="flex items-center gap-2 px-4 py-2 rounded-lg border border-red-300 text-sm font-medium text-red-600 hover:bg-red-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                    {isSubmitting ? (
                        <Spinner />
                    ) : (
                        <TiDeleteOutline className="w-4 h-4 text-red-500" size={20} />
                    )}
                    Reject
                </button>
            </div>
        );
    }

    return (
        <div className="flex items-center gap-3">
            <button
                onClick={onReject}
                disabled={isSubmitting}
                className="flex items-center gap-2 px-4 py-2 rounded-lg border border-gray-300 text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
                {isSubmitting ? <Spinner /> : <AiOutlineClose className="w-4 h-4" />}
                Reject
            </button>
            <button
                onClick={onApprove}
                disabled={isSubmitting}
                className="flex items-center gap-2 px-5 py-2 rounded-lg bg-[#1a3fbf] text-sm font-semibold text-white hover:bg-[#1535a8] disabled:opacity-50 disabled:cursor-not-allowed transition-colors shadow-sm"
            >
                {isSubmitting ? (
                    <Spinner light />
                ) : (
                    <IoCheckmark className="w-4 h-4" />
                )}
                Approve
            </button>
        </div>
    );
}


function Spinner({ light }: { light?: boolean }) {
    return (
        <svg
            className={`w-4 h-4 animate-spin ${light ? "text-white" : "text-gray-500"
                }`}
            viewBox="0 0 24 24"
            fill="none"
        >
            <circle
                className="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                strokeWidth="4"
            />
            <path
                className="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"
            />
        </svg>
    );
}