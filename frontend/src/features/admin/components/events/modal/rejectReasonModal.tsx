import { useState } from "react";

// src/features/admin/components/events/modal/rejectReasonModal.tsx
interface Props {
    isOpen: boolean;
    onClose: () => void;
    onSubmit: (reason: string) => void;
    isPending?: boolean;
}

export default function RejectReasonModal({ isOpen, onClose, onSubmit, isPending }: Props) {
    const [reason, setReason] = useState<string>("");

    if (!isOpen) return null;

    const handleSubmit = () => {
        if (!reason.trim()) return;
        onSubmit(reason);
        setReason("");
    };

    const handleClose = () => {
        setReason("");
        onClose();
    };

    return (
        <div
            className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center"
            onClick={handleClose}
        >
            <div
                className="bg-white p-6 rounded-xl w-full max-w-md shadow-lg"
                onClick={(e) => e.stopPropagation()}
            >
                <h2 className="font-bold text-lg text-gray-900 mb-1">Reject Event</h2>
                <p className="text-sm text-gray-500 mb-4">
                    Please provide a reason for rejection. This will be sent to the owner.
                </p>

                <textarea
                    className="w-full p-3 rounded-lg bg-gray-100 outline-none resize-none text-sm"
                    rows={4}
                    placeholder="Enter rejection reason..."
                    value={reason}
                    onChange={(e) => setReason(e.target.value)}
                />

                <div className="flex justify-end gap-3 mt-4">
                    <button
                        onClick={handleClose}
                        className="px-4 py-2 border border-gray-300 rounded-lg text-sm hover:bg-gray-50 transition-colors"
                    >
                        Cancel
                    </button>
                    <button
                        onClick={handleSubmit}
                        disabled={!reason.trim() || isPending}
                        className="px-4 py-2 bg-red-500 text-white rounded-lg text-sm hover:bg-red-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        {isPending ? "Submitting..." : "Submit"}
                    </button>
                </div>
            </div>
        </div>
    );
}