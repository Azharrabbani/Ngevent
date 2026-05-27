import StatusBadge from "./statusBadge";
import ApprovalButtons from "./approvalButtons";
import BackButton from "./backButton";
import InfoField from "./infoField";
import { useState } from "react";
import { useGetEventByID } from "../../../events/hooks/useGetEventByID";
import { formatTimeAgo, timeRange } from "../../../../utils/timeRange";
import { toDateString } from "../../../../utils/dateConverter";
import MapPicker from "../../../../components/map";
import { PinIcon } from "../../../../components/icon";
import DetailCard from "./detailCard";
import ViewBannerModal from "./modal/bannerModal";
import Modal from "../../../../components/modal";
import RejectReasonModal from "./modal/rejectReasonModal";

interface EventDetailViewProps {
    id?: string;
    onBack?: () => void;
    onApprove?: () => Promise<void> | void;
    onReject?: (reason: string) => Promise<void> | void;
    isSubmitting?: boolean;
}

export default function EventDetailView({
    id,
    onBack,
    onApprove = async () => { },
    onReject = async () => { },
    isSubmitting }: EventDetailViewProps) {
    const [actionDone, setActionDone] = useState<"approved" | "rejected" | null>(
        null
    );
    const [openBanner, setOpenBanner] = useState<boolean>(false);
    const [rejectOpen, setRejectOpen] = useState(false);

    const handleApprove = async () => {
        await onApprove();
        setActionDone("approved");
    };

    const handleReject = async (reason: string) => {
        await onReject(reason);
        setActionDone("rejected");
    };

    const { data, isLoading } = useGetEventByID(id!);

    const start = new Date(Number(data?.start_time) * 1000)
    const end = new Date(Number(data?.end_time) * 1000)

    const date = toDateString(start, "long")
    const eventTimeRange = timeRange(start, end);
    const submittedAt = new Date(Number(data?.created_at) * 1000)

    const isPending = data?.event?.status?.toLowerCase() === "pending";

    if (isLoading) {
        return (
            <div className="fixed inset-0 z-50 flex items-center justify-center bg-gray-50">
                <div className="text-center">
                    <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
                    <p className="text-blue-600 font-medium">Loading event...</p>
                </div>
            </div>
        );
    }

    return (

        <div className="bg-gray-50 p-4 sm:p-6 lg:p-8 font-sans">
            <div className="max-w-4xl mx-auto">
                <BackButton onClick={onBack} />

                {actionDone && (
                    <div
                        className={`mb-4 px-4 py-3 rounded-lg text-sm font-medium flex items-center gap-2 ${actionDone === "approved"
                            ? "bg-green-50 text-green-700 border border-green-200"
                            : "bg-red-50 text-red-700 border border-red-200"
                            }`}
                    >
                        <span>{actionDone === "approved" ? "✓" : "✕"}</span>
                        Event has been{" "}
                        {actionDone === "approved" ? "approved" : "rejected"} successfully.
                    </div>
                )}

                <DetailCard>
                    <div className="w-full h-48 sm:h-64 bg-gray-100 overflow-hidden">
                        {data?.event?.banner ? (
                            <img
                                src={data?.event?.banner}
                                alt={data?.event?.name}
                                className="w-full h-full object-cover hover:opacity-90 transition cursor-pointer"
                                onClick={() => setOpenBanner(true)}
                            />
                        ) : (
                            <div className="w-full h-full flex items-center justify-center bg-gray-100">
                                <span className="text-gray-400 text-sm">No banner available</span>
                            </div>
                        )}
                    </div>

                    {/* Header */}
                    <div className="p-5 sm:p-6 flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4 border-b border-gray-100">
                        <div>
                            <div className="flex items-center gap-3 mb-1">
                                <h1 className="text-xl sm:text-2xl font-bold text-gray-900">
                                    {data?.event.name}
                                </h1>
                                <StatusBadge status={data?.event?.status.toUpperCase() ?? "PENDING"} />
                            </div>
                            <p className="text-sm text-gray-400">
                                ID: {data?.id} &bull; Submitted {formatTimeAgo(submittedAt)}
                            </p>
                        </div>

                        {!actionDone && isPending && (
                            <div className="shrink-0">
                                <ApprovalButtons
                                    variant="outline"
                                    isSubmitting={isSubmitting}
                                    onApprove={handleApprove}
                                    onReject={() => setRejectOpen(true)}
                                />
                            </div>
                        )}
                    </div>

                    <div className="px-5 sm:px-6 py-4 border-b border-gray-100 grid grid-cols-2 sm:grid-cols-4 gap-4">
                        <InfoField label="Date">
                            {date}
                        </InfoField>
                        <InfoField label="Time">
                            {eventTimeRange}
                            <span className="block text-gray-400 text-xs">
                                WIB
                            </span>
                        </InfoField>
                        <InfoField label="City">{data?.event_address?.city}</InfoField>
                        <InfoField label="Country">{data?.event_address?.country}</InfoField>
                    </div>

                    <div className="p-5 sm:p-6 flex flex-col lg:flex-row gap-6">
                        <div className="flex-1">
                            <h2 className="text-base font-bold text-gray-900 mb-3">
                                Event Description
                            </h2>
                            <div
                                className="prose prose-sm text-sm text-gray-600 leading-relaxed space-y-3"
                                dangerouslySetInnerHTML={{ __html: data?.event.description || "" }}
                            />
                        </div>

                        <div className="hidden lg:block w-px bg-gray-100" />
                        <div className="lg:hidden h-px bg-gray-100" />

                        <div className="w-full lg:w-72 shrink-0">
                            <h2 className="text-base font-bold text-gray-900 mb-3 flex items-center gap-2">
                                <PinIcon className="w-4 h-4 text-gray-500" />
                                Location Details
                            </h2>
                            <MapPicker
                                position={[data?.event_address?.coordinates?.lat ?? 0, data?.event_address?.coordinates?.lon ?? 0]}
                                selectedLocation={{
                                    display_name: data?.event_address?.address || "",
                                    lat: String(data?.event_address?.coordinates?.lat),
                                    lon: String(data?.event_address?.coordinates?.lon)
                                }} />
                            <InfoField label="Address" className="mb-3">
                                {data?.event_address?.address}
                            </InfoField>
                            <InfoField label="Detail Address">
                                {data?.event_address?.detail_address}
                            </InfoField>
                        </div>
                    </div>
                </DetailCard>

                <Modal
                    isOpen={openBanner}
                    onClose={() => setOpenBanner(false)}
                >
                    <ViewBannerModal
                        bannerUrl={data?.event?.banner || ""}
                    />
                </Modal>

                <RejectReasonModal
                    isOpen={rejectOpen}
                    onClose={() => setRejectOpen(false)}
                    onSubmit={async (reason) => {
                        setRejectOpen(false);
                        await handleReject(reason);
                    }}
                />
            </div>
        </div>


    )
}
