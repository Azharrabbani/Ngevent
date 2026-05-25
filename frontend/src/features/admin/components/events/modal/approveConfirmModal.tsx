interface Props {
    isOpen: boolean;
    eventName: string;
    onClose: () => void;
    onView: () => void;
    onApprove: () => void;
    isPending?: boolean;
}

export default function ApproveConfirmModal({ isOpen, eventName, onClose, onView, onApprove, isPending }: Props) {
    if (!isOpen) return null;

    return (
        <div
            className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center"
            onClick={onClose}
        >
            <div
                className="bg-white p-6 rounded-xl w-full max-w-md shadow-lg"
                onClick={(e) => e.stopPropagation()}
            >
                <h2 className="font-bold text-lg text-gray-900 mb-1">Approve Event</h2>
                <p className="text-sm text-gray-500 mb-1">
                    You are about to approve:
                </p>
                <p className="text-sm font-semibold text-gray-800 mb-4 bg-gray-50 px-3 py-2 rounded-lg">
                    {eventName}
                </p>
                <p className="text-sm text-gray-500 mb-6">
                    You can review the full event details before approving, or approve directly if you're confident.
                </p>

                <div className="flex justify-end gap-3">
                    <button
                        onClick={onClose}
                        className="px-4 py-2 border border-gray-300 rounded-lg text-sm hover:bg-gray-50 transition-colors"
                    >
                        Cancel
                    </button>
                    <button
                        onClick={onView}
                        className="px-4 py-2 border border-blue-500 text-blue-600 rounded-lg text-sm hover:bg-blue-50 transition-colors"
                    >
                        View Details
                    </button>
                    <button
                        onClick={onApprove}
                        disabled={isPending}
                        className="px-4 py-2 bg-blue-500 text-white rounded-lg text-sm hover:bg-blue-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        {isPending ? "Approving..." : "Approve"}
                    </button>
                </div>
            </div>
        </div>
    );
}