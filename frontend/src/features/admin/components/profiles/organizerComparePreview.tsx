import { useState } from "react";
import { useValidateOrganizerUpdate } from "../../../profile/hooks/organizer/useValidateOrganizer";
import type { OrganizerResponse, OrganizerUpdateResponse } from "../../../profile/types/profileResponse";
import FilePreview from "../../../../components/filePreview";
import { FaRegFilePdf } from "react-icons/fa";

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

    const isPending = requestUpdate?.status === "pending";

    const handleApprove = async () => {
        if (!requestUpdate?.id) return;

        await validateOrganizer(
            {
                id: requestUpdate.id,
                payload: { status: "approved", reason: "" },
            },
            {
                onSuccess: () => onClose?.(),
            }
        );
    };

    const handleReject = async () => {
        if (!requestUpdate?.id) return;

        await validateOrganizer(
            {
                id: requestUpdate.id,
                payload: { status: "rejected", reason },
            },
            {
                onSuccess: () => onClose?.(),
            }
        );
    };

    if (profileLoading || updateLoading) {
        return (
            <div className="grid grid-cols-1 xl:grid-cols-2 gap-6 mt-4 animate-pulse">
                <div className="h-[500px] rounded-2xl bg-gray-100 border border-gray-200" />
                <div className="h-[500px] rounded-2xl bg-blue-50 border border-blue-100" />
            </div>
        );
    }

    return (
        <>
            <div className="grid grid-cols-1 xl:grid-cols-2 gap-6 mt-4">
                {/* Current Data */}
                <div className="bg-white border border-gray-200 rounded-2xl shadow-sm overflow-hidden">
                    <div className="px-5 py-4 border-b bg-gray-50">
                        <h3 className="font-bold text-lg text-gray-800">
                            Current Data
                        </h3>
                        <p className="text-sm text-gray-500 mt-1">
                            Event owner current profile information
                        </p>
                    </div>

                    <div className="p-5 space-y-4 text-sm">
                        <InfoItem label="Email" value={current.email} />
                        <InfoItem
                            label="Secondary Email"
                            value={current.social_media?.email || "-"}
                        />
                        <InfoItem label="Phone" value={current.phone_number} />
                        <InfoItem label="Country" value={current.country} />
                        <InfoItem label="Address" value={current.address} />

                        <div>
                            <p className="font-semibold mb-2">Description</p>
                            <div
                                className="prose prose-sm max-w-none bg-gray-50 border border-gray-200 rounded-xl p-4"
                                dangerouslySetInnerHTML={{
                                    __html: current.company_detail?.description || "-",
                                }}
                            />
                        </div>

                        {/* Current legal documents */}
                        <div>
                            <p className="font-semibold mb-2">Legal Documents</p>
                            <div className="grid grid-cols-2 gap-3">
                                <DocumentCard
                                    label="NPWP"
                                    number={current.company_detail?.npwp}
                                    fileUrl={current.company_detail?.npwp_file}
                                    onPreview={setPreviewFile}
                                    variant="default"
                                />
                                <DocumentCard
                                    label="NIB"
                                    number={current.company_detail?.nib}
                                    fileUrl={current.company_detail?.nib_file}
                                    onPreview={setPreviewFile}
                                    variant="default"
                                />
                            </div>
                        </div>
                    </div>
                </div>

                {/* Requested Update */}
                <div className="bg-gradient-to-br from-blue-50 to-white border border-blue-200 rounded-2xl shadow-sm overflow-hidden">
                    <div className="px-5 py-4 border-b border-blue-100">
                        <div className="flex items-center justify-between">
                            <div>
                                <h3 className="font-bold text-lg text-blue-800">
                                    Requested Update
                                </h3>
                                <p className="text-sm text-blue-500 mt-1">
                                    New event owner information submitted
                                </p>
                            </div>
                            <span className="px-3 py-1 rounded-full bg-yellow-100 text-yellow-700 text-xs font-semibold">
                                Pending Review
                            </span>
                        </div>
                    </div>

                    <div className="p-5 space-y-4 text-sm">
                        <InfoItem label="Name" value={requestUpdate.name} />
                        <InfoItem label="Secondary Email" value={requestUpdate.email} />
                        <InfoItem label="Phone" value={requestUpdate.phone_number} />
                        <InfoItem label="Country" value={requestUpdate.country} />
                        <InfoItem label="Address" value={requestUpdate.address} />

                        <div>
                            <p className="font-semibold mb-2">Description</p>
                            <div
                                className="prose prose-sm max-w-none bg-white border border-blue-100 rounded-xl p-4"
                                dangerouslySetInnerHTML={{
                                    __html: requestUpdate.description || "-",
                                }}
                            />
                        </div>

                        {/* Requested legal documents */}
                        <div>
                            <p className="font-semibold mb-2">Legal Documents</p>
                            <div className="grid grid-cols-2 gap-3">
                                <DocumentCard
                                    label="NPWP"
                                    number={requestUpdate.npwp_number}
                                    fileUrl={requestUpdate.npwp_document}
                                    onPreview={setPreviewFile}
                                    variant="update"
                                    fallbackFileUrl={current.company_detail?.npwp_file}
                                    fallbackNumber={current.company_detail?.npwp}
                                />
                                <DocumentCard
                                    label="NIB"
                                    number={requestUpdate.nib_number}
                                    fileUrl={requestUpdate.nib_document}
                                    onPreview={setPreviewFile}
                                    variant="update"
                                    fallbackFileUrl={current.company_detail?.nib_file}
                                    fallbackNumber={current.company_detail?.nib}
                                />
                            </div>
                        </div>
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
                        onClick={() => setRejectModalOpen(true)}
                        className="px-5 py-2.5 rounded-xl border border-red-200 text-red-600 hover:bg-red-50 transition font-medium"
                    >
                        Reject
                    </button>

                    <button
                        onClick={handleApprove}
                        disabled={updatePending}
                        className="px-5 py-2.5 rounded-xl bg-blue-600 text-white hover:bg-blue-700 transition font-medium disabled:opacity-60"
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
                        <h2 className="font-bold text-lg mb-3">Reject Event Owner</h2>

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
                                className="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 disabled:opacity-60"
                            >
                                {updatePending ? "Submitting..." : "Submit"}
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </>
    );
}

function InfoItem({ label, value }: { label: string; value: string }) {
    return (
        <div>
            <p className="text-xs font-medium text-gray-500 mb-1">{label}</p>
            <div className="bg-gray-50 border border-gray-200 rounded-xl px-3 py-2 text-gray-800">
                {value || "-"}
            </div>
        </div>
    );
}

interface DocumentCardProps {
    label: string;
    number?: string;
    fileUrl?: string;
    onPreview: (url: string) => void;
    variant: "default" | "update";
    fallbackFileUrl?: string;
    fallbackNumber?: string;
}

function DocumentCard({
    label,
    number,
    fileUrl,
    onPreview,
    variant,
    fallbackFileUrl,
    fallbackNumber,
}: DocumentCardProps) {
    const isUpdate = variant === "update";
    const resolvedUrl = fileUrl || fallbackFileUrl;
    const resolvedNumber = number || fallbackNumber;
    const isUnchanged = isUpdate && !fileUrl && !!fallbackFileUrl;

    const borderClass = isUpdate ? "border-blue-100 bg-white" : "border-gray-200 bg-gray-50";
    const buttonClass = isUpdate
        ? "bg-blue-600 hover:bg-blue-700 text-white"
        : "bg-red-500  hover:bg-red-600  text-white";

    return (
        <div className={`border rounded-xl p-3 ${borderClass}`}>
            <p className="text-xs text-gray-500 mb-1">{label} Number</p>
            <p className="font-semibold text-gray-800 text-sm truncate">
                {resolvedNumber || "-"}
            </p>

            {isUnchanged && (
                <p className="text-xs text-gray-400 mt-1 italic">No new file submitted</p>
            )}

            {resolvedUrl ? (
                <button
                    onClick={() => onPreview(resolvedUrl)}
                    className={`mt-3 flex items-center gap-2 px-3 py-1.5 rounded-lg text-xs font-medium transition ${buttonClass}`}
                >
                    <FaRegFilePdf />
                    Preview {label}
                </button>
            ) : (
                <p className="text-xs text-gray-400 mt-3 italic">No file available</p>
            )}
        </div>
    );
}