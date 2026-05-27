import { FaInstagram, FaRegFilePdf } from "react-icons/fa";
import { useState } from "react";
import { useValidateOrganizerUpdate } from "../../../profile/hooks/organizer/useValidateOrganizer";
import type { OrganizerResponse, OrganizerUpdateResponse } from "../../../profile/types/profileResponse";
import FilePreview from "../../../../components/filePreview";

interface Props {
    requestUpdate: OrganizerUpdateResponse;
    current: OrganizerResponse;
    profileLoading: boolean | undefined;
    updateLoading: boolean | undefined;
    onClose?: () => void;
};

export default function OrganizerComparePreview({
    requestUpdate,
    current,
    profileLoading,
    updateLoading,
    onClose
}: Props) {
    const [previewFile, setPreviewFile] = useState<string | null>(null);
    const [rejectModalOpen, setRejectModalOpen] = useState(false);
    const [reason, setReason] = useState("");

    const { mutateAsync: validateOrganizer, isPending: updatePending } = useValidateOrganizerUpdate();

    const isPending = current?.status?.status === "pending";

    const handleApprove = async () => {
        if (!requestUpdate?.id) return;

        await validateOrganizer({
            id: requestUpdate?.id,
            payload: {
                status: "approved",
                reason: "",
            },
        });

        onClose?.();
    };

    const handleReject = async () => {
        if (!requestUpdate?.id) return;

        await validateOrganizer({
            id: requestUpdate?.id,
            payload: {
                status: "rejected",
                reason,
            },
        });

        setRejectModalOpen(false);
        onClose?.();
    };

    return (
        <>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mt-2">
                <div className="border-2 border-black rounded-xl p-4 shadow-[3px_3px_0px_#000]">
                    {profileLoading && <h3>Loading...</h3>}
                    <h3 className="font-bold mb-3">Current Data</h3>
                    <div className="space-y-2 text-sm">
                        <p><b>Email:</b> {current.email}</p>
                        <p><b>Secondary Email:</b> {current.social_media.email ? current.social_media.email : "-"}</p>
                        <p><b>Phone:</b> {current.phone_number}</p>
                        <p><b>Country:</b> {current.country}</p>
                        <p><b>Address:</b> {current.address}</p>
                        <p><b>Description:</b></p>
                        <div
                            className="prose prose-sm max-w-none bg-gray-50 border border-gray-200 rounded-lg p-3 text-sm"
                            dangerouslySetInnerHTML={{ __html: current.company_detail?.description }}
                        />
                        <p><b>NPWP:</b> {current.company_detail?.npwp}</p>
                        <p className="flex items-center gap-2">
                            <b>File:</b>
                            <FaRegFilePdf
                                className="cursor-pointer text-red-500 hover:scale-110 transition"
                                onClick={() => setPreviewFile(current.company_detail?.npwp_file || null)}
                            />
                        </p>
                        <p><b>NIB:</b> {current.company_detail?.nib}</p>
                        <p className="flex items-center gap-2">
                            <b>File:</b>
                            <FaRegFilePdf
                                className="cursor-pointer text-red-500 hover:scale-110 transition"
                                onClick={() => setPreviewFile(current.company_detail?.nib_file || null)}
                            />
                        </p>
                        <p className="flex items-center gap-2">
                            <b>Instagram:</b>
                            {current.social_media?.instagram ? (
                                <FaInstagram
                                    className="cursor-pointer text-pink-500 hover:scale-110 transition"
                                    onClick={() =>
                                        window.open(current.social_media!.instagram, "_blank")
                                    }
                                />
                            ) : "-"}
                        </p>
                    </div>
                </div>
                <div className="border-2 border-dashed border-gray-400 rounded-xl p-4 bg-gray-50">
                    {updateLoading && <h3>Loading...</h3>}
                    <h3 className="font-bold mb-3 text-gray-600">
                        Requested Update
                    </h3>
                    <div className="space-y-2 text-sm text-gray-700">
                        <p><b>Name:</b> {requestUpdate.name}</p>
                        <p><b>Secondary Email:</b> {requestUpdate.email}</p>
                        <p><b>Phone:</b> {requestUpdate.phone_number}</p>
                        <p><b>Country:</b> {requestUpdate.country}</p>
                        <p><b>Address:</b> {requestUpdate.address}</p>
                        <p><b>Description:</b></p>
                        <div
                            className="prose prose-sm max-w-none bg-gray-50 border border-dashed border-gray-300 rounded-lg p-3 text-sm"
                            dangerouslySetInnerHTML={{ __html: requestUpdate.description }}
                        />
                        <p><b>NPWP:</b> {requestUpdate.npwp_number}</p>
                        {requestUpdate.npwp_document &&
                            <p className="flex items-center gap-2">
                                <b>NPWP File:</b>
                                <FaRegFilePdf
                                    className="cursor-pointer text-red-500 hover:scale-110 transition"
                                    onClick={() => setPreviewFile(requestUpdate.npwp_document || null)}
                                />
                            </p>
                        }
                        <p><b>NIB:</b> {requestUpdate.nib_number}</p>
                        {requestUpdate.nib_document &&
                            <p className="flex items-center gap-2">
                                <b>NIB File:</b>
                                <FaRegFilePdf
                                    className="cursor-pointer text-red-500 hover:scale-110 transition"
                                    onClick={() => setPreviewFile(requestUpdate.nib_document || null)}
                                />
                            </p>
                        }
                    </div>
                </div>
            </div>

            {previewFile && (
                <FilePreview
                    setPreview={() => setPreviewFile(null)}
                    preview={previewFile}
                />
            )}

            {isPending && (
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
                        disabled={updatePending}
                    >
                        {updatePending ? "Processing..." : "Approve"}
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
                                disabled={!reason || updatePending}
                                className="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600"
                            >
                                {updatePending ? "Submitting..." : "Submit"}
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </>


    )
}