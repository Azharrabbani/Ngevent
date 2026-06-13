import { FaRegFilePdf } from "react-icons/fa";
import { HiOutlineMail, HiOutlinePhone, HiOutlineLocationMarker } from "react-icons/hi";
import { useState } from "react";
import type { OrganizerResponse } from "../../../profile/types/profileResponse";
import { useApproveOrganizer } from "../../../profile/hooks/organizer/useApproveOrganizer";
import { useRejectOrganizer } from "../../../profile/hooks/organizer/useRejectOrganzier";
import FilePreview from "../../../../components/filePreview";
import { toDateString } from "../../../../utils/dateConverter";
import { InstagramIcon } from "../../../../components/icon";

interface Props {
    loading: boolean | undefined;
    profile: OrganizerResponse;
    hasUpdate?: boolean;
    onClose?: () => void;
}

export default function OrganizerPreview({
    loading,
    profile,
    hasUpdate,
    onClose,
}: Props) {
    const [previewFile, setPreviewFile] = useState<string | null>(null);
    const [rejectModalOpen, setRejectModalOpen] = useState(false);
    const [reason, setReason] = useState("");

    const { mutateAsync: approval, isPending: approvalPending } =
        useApproveOrganizer(profile.id);

    const { mutateAsync: reject, isPending: rejectPending } =
        useRejectOrganizer(profile.id);

    const reviewedAt = new Date(Number(profile?.status?.reviewed_at) * 1000)
    const reviewDate = toDateString(reviewedAt, "long")

    const isPending = profile?.status?.status === "pending";
    const isRejected = profile?.status?.status === "rejected";

    const showAction = isPending && !hasUpdate;

    const handleApprove = async () => {
        if (!profile?.id) return;

        await approval(profile.id);

        onClose?.();
    };

    const handleReject = async () => {
        if (!profile?.id) return;

        await reject({
            id: profile.id,
            payload: {
                reason,
            },
        });

        setRejectModalOpen(false);
        onClose?.();
    };

    return (
        <>
            <div className="mt-5 space-y-6 text-sm sm:text-base break-words">
                {loading && (
                    <div className="animate-pulse text-gray-500">
                        Loading organizer data...
                    </div>
                )}

                {/* Header */}
                <div className="flex items-start gap-4 border border-gray-200 rounded-2xl p-5 bg-gradient-to-br from-white to-gray-50 shadow-sm">
                    <div className="w-16 h-16 rounded-2xl bg-blue-100 flex items-center justify-center text-2xl font-bold text-blue-600">
                        {profile?.name?.charAt(0) || "O"}
                    </div>

                    <div className="flex-1">
                        <div className="flex flex-wrap items-center gap-3">
                            <h2 className="text-xl font-bold text-gray-800">
                                {profile?.name}
                            </h2>

                            <span
                                className={`px-3 py-1 rounded-full text-xs font-semibold capitalize
                                ${(profile?.status?.status === "approved")
                                        ? "bg-green-100 text-green-700"
                                        : (profile?.status?.status === "rejected" || profile?.status?.status === "deactivated")
                                            ? "bg-red-100 text-red-700"
                                            : "bg-yellow-100 text-yellow-700"
                                    }`}
                            >
                                {profile?.status?.status}
                            </span>
                        </div>

                        <div className="mt-3 space-y-2 text-gray-600">
                            <p className="flex items-center gap-2">
                                <HiOutlineMail className="text-gray-500" />
                                {profile?.email}
                            </p>

                            <p className="flex items-center gap-2">
                                <HiOutlinePhone className="text-gray-500" />
                                {profile?.phone_number}
                            </p>

                            <p className="flex items-center gap-2">
                                <HiOutlineLocationMarker className="text-gray-500" />
                                {profile?.country}
                            </p>
                        </div>
                    </div>
                </div>

                <div className="border border-gray-200 rounded-2xl p-5 bg-white shadow-sm">
                    <h3 className="text-lg font-bold text-gray-800 mb-4">
                        Company Detail
                    </h3>

                    <div className="space-y-4">
                        <div>
                            <p className="font-semibold text-gray-700 mb-2">
                                Address
                            </p>

                            <div className="bg-gray-50 border border-gray-200 rounded-xl p-3 text-gray-700">
                                {profile?.address || "-"}
                            </div>
                        </div>

                        <div>
                            <p className="font-semibold text-gray-700 mb-2">
                                Description
                            </p>

                            <div
                                className="prose prose-sm max-w-none bg-gray-50 border border-gray-200 rounded-xl p-4"
                                dangerouslySetInnerHTML={{
                                    __html:
                                        profile?.company_detail?.description ||
                                        "<p>-</p>",
                                }}
                            />
                        </div>
                    </div>
                </div>

                <div className="border border-gray-200 rounded-2xl p-5 bg-white shadow-sm">
                    <h3 className="text-lg font-bold text-gray-800 mb-4">
                        Legal Documents
                    </h3>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        {/* NPWP */}
                        <div className="border border-gray-200 rounded-xl p-4 bg-gray-50">
                            <p className="text-sm text-gray-500 mb-1">
                                NPWP Number
                            </p>

                            <p className="font-semibold text-gray-800">
                                {profile?.company_detail?.npwp || "-"}
                            </p>

                            <button
                                onClick={() =>
                                    setPreviewFile(
                                        profile?.company_detail?.npwp_file || null
                                    )
                                }
                                className="mt-4 flex items-center gap-2 px-4 py-2 rounded-lg bg-red-500 text-white hover:bg-red-600 transition"
                            >
                                <FaRegFilePdf />
                                Preview NPWP
                            </button>
                        </div>

                        {/* NIB */}
                        <div className="border border-gray-200 rounded-xl p-4 bg-gray-50">
                            <p className="text-sm text-gray-500 mb-1">
                                NIB Number
                            </p>

                            <p className="font-semibold text-gray-800">
                                {profile?.company_detail?.nib || "-"}
                            </p>

                            <button
                                onClick={() =>
                                    setPreviewFile(
                                        profile?.company_detail?.nib_file || null
                                    )
                                }
                                className="mt-4 flex items-center gap-2 px-4 py-2 rounded-lg bg-red-500 text-white hover:bg-red-600 transition"
                            >
                                <FaRegFilePdf />
                                Preview NIB
                            </button>
                        </div>
                    </div>
                </div>

                <div className="border border-gray-200 rounded-2xl p-5 bg-white shadow-sm">
                    <h3 className="text-lg font-bold text-gray-800 mb-4">
                        Social Media
                    </h3>

                    <div className="space-y-4">
                        <div>
                            <p className="text-sm text-gray-500 mb-1">
                                Secondary Email
                            </p>

                            <div className="bg-gray-50 border border-gray-200 rounded-xl p-3 text-gray-700">
                                {profile?.social_media?.email || "-"}
                            </div>
                        </div>

                        <div>
                            <p className="text-sm text-gray-500 mb-2">
                                Instagram
                            </p>

                            {profile?.social_media?.instagram ? (
                                <button
                                    onClick={() =>
                                        window.open(
                                            profile.social_media.instagram,
                                            "_blank"
                                        )
                                    }
                                    className="flex items-center gap-2 px-4 py-2 rounded-xl bg-gradient-to-r from-pink-500 to-purple-500 text-white hover:opacity-90 transition"
                                >
                                    <InstagramIcon className="text-lg" />
                                    Visit Instagram
                                </button>
                            ) : (
                                <div className="bg-gray-50 border border-gray-200 rounded-xl p-3 text-gray-700">
                                    -
                                </div>
                            )}
                        </div>
                    </div>
                </div>

                {isRejected && (
                    <div className="border border-red-200 rounded-2xl p-5 bg-red-50 shadow-sm">
                        <h3 className="text-lg font-bold text-red-700 mb-4">
                            Rejection Review
                        </h3>

                        <div className="space-y-4">
                            <div>
                                <p className="text-sm font-medium text-gray-500 mb-1">
                                    Reviewed By
                                </p>

                                <div className="bg-white border border-red-100 rounded-xl p-3 text-gray-800">
                                    {profile?.status.reviewed_by || "-"}
                                </div>
                            </div>

                            <div>
                                <p className="text-sm font-medium text-gray-500 mb-1">
                                    Reason
                                </p>

                                <div className="bg-white border border-red-100 rounded-xl p-3 text-gray-800">
                                    {profile?.status.rejected_reason || "-"}
                                </div>
                            </div>

                            <div>
                                <p className="text-sm font-medium text-gray-500 mb-1">
                                    Reviewed At
                                </p>

                                <div className="bg-white border border-red-100 rounded-xl p-3 text-gray-800">
                                    {reviewDate ?? "-"}
                                </div>
                            </div>
                        </div>
                    </div>
                )}

                {previewFile && (
                    <FilePreview
                        setPreview={() => setPreviewFile(null)}
                        preview={previewFile}
                    />
                )}
            </div>

            {showAction && (
                <div className="flex justify-end gap-3 mt-8">
                    <button
                        className="px-5 py-2.5 border border-gray-300 rounded-xl hover:bg-gray-100 transition font-medium"
                        onClick={() => setRejectModalOpen(true)}
                    >
                        Reject
                    </button>

                    <button
                        className="px-5 py-2.5 bg-blue-600 text-white rounded-xl hover:bg-blue-700 transition font-medium disabled:opacity-60"
                        onClick={handleApprove}
                        disabled={approvalPending}
                    >
                        {approvalPending ? "Processing..." : "Approve"}
                    </button>
                </div>
            )}

            {rejectModalOpen && (
                <div
                    className="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4"
                    onClick={() => setRejectModalOpen(false)}
                >
                    <div
                        className="bg-white w-full max-w-md rounded-3xl p-6 shadow-2xl"
                        onClick={(e) => e.stopPropagation()}
                    >
                        <h2 className="text-xl font-bold text-gray-800 mb-2">
                            Reject Organizer
                        </h2>

                        <p className="text-sm text-gray-500 mb-4">
                            Please provide the reason for rejecting this organizer.
                        </p>

                        <textarea
                            className="w-full p-4 rounded-2xl bg-gray-100 border border-gray-200 outline-none resize-none focus:ring-2 focus:ring-red-400"
                            rows={5}
                            placeholder="Enter rejection reason..."
                            value={reason}
                            onChange={(e) => setReason(e.target.value)}
                        />

                        <div className="flex justify-end gap-3 mt-5">
                            <button
                                onClick={() => setRejectModalOpen(false)}
                                className="px-4 py-2 rounded-xl border border-gray-300 hover:bg-gray-100 transition"
                            >
                                Cancel
                            </button>

                            <button
                                onClick={handleReject}
                                disabled={!reason || rejectPending}
                                className="px-4 py-2 rounded-xl bg-red-500 text-white hover:bg-red-600 transition disabled:opacity-60"
                            >
                                {rejectPending ? "Submitting..." : "Submit"}
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </>
    );
}