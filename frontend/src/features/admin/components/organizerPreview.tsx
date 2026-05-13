import { FaInstagram, FaRegFilePdf } from "react-icons/fa";
import type { OrganizerResponse } from "../../profile/types/profileResponse";
import { useState } from "react";
import FilePreview from "../../../components/filePreview";
import { useApproveOrganizer } from "../../profile/hooks/organizer/useApproveOrganizer";
import { useRejectOrganizer } from "../../profile/hooks/organizer/useRejectOrganzier";

interface Props {
    loading: boolean | undefined;
    profile: OrganizerResponse,
    hasUpdate?: boolean;
    onClose?: () => void;
}

export default function OrganizerPreview({ loading, profile, hasUpdate, onClose }: Props) {
    const [previewFile, setPreviewFile] = useState<string | null>(null);
    const [rejectModalOpen, setRejectModalOpen] = useState(false);
    const [reason, setReason] = useState("");

    const { mutateAsync: approval, isPending: approvalPending } = useApproveOrganizer(profile.id);

    const { mutateAsync: reject, isPending: rejectPending } = useRejectOrganizer(profile.id);

    const isPending = profile?.status?.status === "pending";
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
            <div className="mt-5 space-y-2 text-sm sm:text-base wrap-break-word">
                {loading && <h3>Loading...</h3>}
                <p><span className="font-semibold">Email:</span> {profile?.email}</p>
                <p><span className="font-semibold">Phone:</span> {profile?.phone_number}</p>
                <p><span className="font-semibold">Country:</span> {profile?.country}</p>
                <p><span className="font-semibold">Address:</span> {profile?.address}</p>
                <p><span className="font-semibold">Description:</span></p>
                <div
                    className="prose prose-sm max-w-none bg-gray-50 border border-gray-200 rounded-lg p-3 text-sm"
                    dangerouslySetInnerHTML={{ __html: profile?.company_detail?.description }}
                />
                <p><span className="font-semibold">NPWP Number:</span> {profile?.company_detail?.npwp}</p>
                <p className="flex gap-2 items-center"><span className="font-semibold">NPWP:</span>
                    <FaRegFilePdf
                        className="cursor-pointer text-2xl text-red-500 hover:scale-110 transition"
                        onClick={() => setPreviewFile(profile?.company_detail?.npwp_file || null)} />
                </p>
                <p><span className="font-semibold">NIB Number:</span> {profile?.company_detail?.nib}</p>
                <p className="flex gap-2 items-center"><span className="font-semibold">NIB:</span>
                    <FaRegFilePdf
                        className="cursor-pointer text-2xl text-red-500 hover:scale-110 transition"
                        onClick={() => setPreviewFile(profile?.company_detail?.nib_file || null)} />
                </p>
                <p><span className="font-semibold">Secondary Email:</span> {profile?.social_media?.email ? profile.social_media?.email : "-"}</p>
                <p className="flex gap-2 items-center"><span className="font-semibold">Instagram:</span>
                    {profile?.social_media?.instagram ?
                        <FaInstagram
                            className="cursor-pointer text-2xl text-pink-500 hover:scale-110 transition"
                            onClick={() => window.open(profile.social_media.instagram, "_blank")}
                        /> :
                        "-"}
                </p>

                {previewFile && (
                    <FilePreview
                        setPreview={() => setPreviewFile(null)}
                        preview={previewFile}
                    />
                )}
            </div>

            {showAction && (
                <div className="flex justify-end gap-3 mt-6">
                    <button
                        className="px-4 py-2 border-2 border-black rounded-lg hover:bg-gray-200 transition"
                        onClick={() => setRejectModalOpen(true)}
                    >
                        Reject
                    </button>
                    <button
                        className="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition"
                        onClick={handleApprove}
                        disabled={approvalPending}
                    >
                        {approvalPending ? "Processing..." : "Approve"}
                    </button>
                </div>
            )}

            {rejectModalOpen && (
                <div
                    className="fixed inset-0 bg-black/50 z-60 flex items-center justify-center"
                    onClick={() => setRejectModalOpen(false)}
                >
                    <div
                        className="bg-white p-6 rounded-xl w-full max-w-md border-2 border-black shadow-[6px_6px_0px_black]"
                        onClick={(e) => e.stopPropagation()}
                    >
                        <h2 className="font-bold text-lg mb-3">Reject Reason</h2>

                        <textarea
                            className="w-full p-3 rounded-lg bg-gray-100 outline-none resize-none"
                            rows={4}
                            placeholder="Enter reason..."
                            value={reason}
                            onChange={(e) => setReason(e.target.value)}
                        />

                        <div className="flex justify-end gap-3 mt-4">
                            <button
                                onClick={() => setRejectModalOpen(false)}
                                className="px-4 py-2 border-2 border-black rounded-lg hover:bg-gray-200"
                            >
                                Cancel
                            </button>

                            <button
                                onClick={handleReject}
                                disabled={!reason || rejectPending}
                                className="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600"
                            >
                                {rejectPending ? "Submitting..." : "Submit"}
                            </button>
                        </div>
                    </div>
                </div>
            )}

        </>
    )
}