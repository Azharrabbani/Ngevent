import { useState } from "react";
import BackButton from "./backButton";
import StatusBadge from "./badge/statusBadge";
import ApprovalButtons from "./approvalButtons";
import DiffRow from "./diffRow";
import { FiEdit } from "react-icons/fi";
import { CalenderIcon, ClockIcon } from "../../../../components/icon";
import { useGetEventByID } from "../../../events/hooks/useGetEventByID";
import { toDateString } from "../../../../utils/dateConverter";
import { formatTimeAgo, timeRange } from "../../../../utils/timeRange";
import { useGetUpdateByEventID } from "../../../events/hooks/useGetUpdateByEventID";
import MapPicker from "../../../../components/map";
import { useParams } from "react-router-dom";
import DetailCard from "./detailCard";
import ViewBannerModal from "./modal/bannerModal";
import BannerPreview from "./modal/bannerPreview";
import Modal from "../../../../components/modal";
import RejectReasonModal from "./modal/rejectReasonModal";
import InfoField from "./infoField";


interface UpdateRequestViewProps {
    status: string | undefined
    onBack?: () => void;
    onApprove?: () => Promise<void> | void;
    onReject?: (reason: string) => Promise<void> | void;
    isSubmitting?: boolean;
}

const isDifferent = (a: unknown, b: unknown) =>
    JSON.stringify(a) !== JSON.stringify(b);

export default function UpdateRequestView({
    onBack,
    status,
    onApprove = async () => { },
    onReject = async () => { },
    isSubmitting }: UpdateRequestViewProps) {
    const { id } = useParams<{ id: string }>();
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
    const { data: currentData, isLoading: currentIsLoading } = useGetEventByID(id!);

    const start = new Date(Number(currentData?.start_time) * 1000)
    const end = new Date(Number(currentData?.end_time) * 1000)

    const currentDate = toDateString(start, "long")
    const currentEventTimeRange = timeRange(start, end);

    const { data: updateData, isLoading: updateIsLoading } = useGetUpdateByEventID(id!, status!);

    const startDate = new Date(Number(updateData?.updated_details?.start_time) * 1000)
    const reviewDate = new Date(Number(updateData?.updated_details?.reviewed_at) * 1000)
    const endDate = new Date(Number(updateData?.updated_details?.end_time) * 1000)

    const updateDate = toDateString(startDate, "long")
    const reviewDateString = toDateString(reviewDate, "long")
    const updateEventTimeRange = timeRange(startDate, endDate);

    const requestedAt = new Date(Number(updateData?.created_at) * 1000)

    const isPending = updateData?.updated_details?.status?.toLowerCase() === "pending";

    if (currentIsLoading || updateIsLoading) {
        return (
            <div className="fixed inset-0 z-50 flex items-center justify-center bg-gray-50">
                <div className="text-center">
                    <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
                    <p className="text-blue-600 font-medium">Loading request...</p>
                </div>
            </div>
        );
    }

    return (
        <div className="bg-gray-50 font-sans">
            <div className="max-w-5xl mx-auto">
                <BackButton onClick={onBack} />

                {/* Header */}
                <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-6">
                    <div>
                        <div className="flex flex-wrap items-center gap-2 mb-1">
                            <StatusBadge status={updateData?.updated_details?.status.toUpperCase() || "PENDING"} />
                            <span className="text-sm text-gray-500">
                                Req ID: {updateData?.id}
                            </span>
                        </div>
                        <h1 className="text-2xl sm:text-3xl font-bold text-gray-900">
                            {currentData?.event?.name}
                        </h1>
                        <p className="text-sm text-gray-400 mt-0.5">
                            Requested by: {currentData?.eo_profile?.name} &bull;{" "}
                            {formatTimeAgo(requestedAt)}
                        </p>
                    </div>
                    {!actionDone && isPending && (
                        <div className="shrink-0">
                            <ApprovalButtons
                                isSubmitting={isSubmitting}
                                onApprove={handleApprove}
                                onReject={() => setRejectOpen(true)}
                            />
                        </div>
                    )}
                </div>

                {actionDone && (
                    <div
                        className={`mb-5 px-4 py-3 rounded-lg text-sm font-medium flex items-center gap-2 ${actionDone === "approved"
                            ? "bg-green-50 text-green-700 border border-green-200"
                            : "bg-red-50 text-red-700 border border-red-200"
                            }`}
                    >
                        <span>{actionDone === "approved" ? "✓" : "✕"}</span>
                        Request has been{" "}
                        {actionDone === "approved" ? "approved" : "rejected"} successfully.
                    </div>
                )}

                <hr className="border-gray-200 mb-6" />

                <DetailCard>
                    <div className="w-full h-48 sm:h-56 bg-gray-100 overflow-hidden border-b border-gray-200">
                        {currentData?.event?.banner ? (
                            <img
                                onClick={() => setOpenBanner(true)}
                                src={currentData?.event?.banner}
                                alt={currentData?.event?.name}
                                className="w-full h-full object-cover cursor-pointer hover:opacity-90 transition"
                            />
                        ) : (
                            <div className="w-full h-full flex items-center justify-center bg-gray-100">
                                <span className="text-gray-400 text-sm">No banner available</span>
                            </div>
                        )}
                    </div>

                    {!isPending && (
                        <div className="px-5 sm:px-6 py-4 border-b border-gray-100 grid grid-cols-1 md:grid-cols-3 gap-4">
                            <InfoField label="Reviewed By">
                                {updateData?.updated_details?.reviewed_by?.email ?? "-"}
                            </InfoField>
                            <InfoField label="Rejected Reason">
                                {updateData?.updated_details?.rejected_reason ?? "-"}
                            </InfoField>
                            <InfoField label="Reviewed at">
                                {reviewDateString ?? "-"}
                            </InfoField>
                        </div>
                    )}

                    {/* Card headers */}
                    <div className="grid grid-cols-1 sm:grid-cols-2 border-b border-gray-200">
                        <div className="px-4 py-3 sm:border-r border-gray-200 flex items-center gap-2">
                            <ClockIcon className="w-4 h-4 text-gray-400" />
                            <span className="text-sm font-semibold text-gray-700">
                                Current Data
                            </span>
                        </div>
                        <div className="px-4 py-3 flex items-center gap-2">
                            <FiEdit className="w-4 h-4 text-gray-400" />
                            <span className="text-sm font-semibold text-gray-700">
                                Proposed Updates
                            </span>
                        </div>
                    </div>

                    {updateData?.updated_details?.banner?.trim() && (
                        <DiffRow
                            label="Banner"
                            currentContent={
                                <BannerPreview
                                    imageUrl={currentData?.event?.banner || ""}
                                    title={currentData?.event?.name || ""}
                                />
                            }
                            proposedContent={
                                <BannerPreview
                                    imageUrl={updateData?.updated_details?.banner}
                                    title={updateData?.event_title || ""}
                                />
                            }
                        />
                    )}

                    <DiffRow
                        label="Event Name"
                        currentContent={currentData?.event?.name}
                        proposedContent={updateData?.event_title}
                        hasChange={isDifferent(currentData?.event?.name, updateData?.event_title)}
                    />

                    <DiffRow
                        label="Date & Time"
                        currentContent={
                            <DateTimeBlock
                                date={currentDate}
                                timeRange={currentEventTimeRange}
                            />
                        }
                        proposedContent={
                            <DateTimeBlock
                                date={updateDate}
                                timeRange={updateEventTimeRange}
                            />
                        }
                        hasChange={
                            isDifferent(currentDate, updateDate) ||
                            isDifferent(currentEventTimeRange, updateEventTimeRange)
                        }
                    />

                    <DiffRow
                        label="Description"
                        currentContent={
                            <DescriptionBlock
                                description={currentData?.event?.description}
                            />
                        }
                        proposedContent={
                            <DescriptionBlock
                                description={updateData?.updated_details?.description}
                            />
                        }
                        hasChange={
                            isDifferent(currentData?.event?.description, updateData?.updated_details?.description)
                        }
                    />

                    <DiffRow
                        label="Location"
                        currentContent={<LocationBlock loc={{
                            address: currentData?.event_address?.address,
                            city: currentData?.event_address?.city,
                            country: currentData?.event_address?.country,
                            detailAddress: currentData?.event_address?.detail_address,
                        }} />}
                        proposedContent={<LocationBlock loc={{
                            address: updateData?.updated_address?.address,
                            city: updateData?.updated_address?.city,
                            country: updateData?.updated_address?.country,
                            detailAddress: updateData?.updated_address?.detail_address,
                        }} />}
                        hasChange={
                            isDifferent(currentData?.event_address?.address, updateData?.updated_address?.address) ||
                            isDifferent(currentData?.event_address?.city, updateData?.updated_address?.city) ||
                            isDifferent(currentData?.event_address?.country, updateData?.updated_address?.country) ||
                            isDifferent(currentData?.event_address?.detail_address, updateData?.updated_address?.detail_address)
                        }
                    />

                    <DiffRow
                        label="Map Marker"
                        currentContent={<MapPicker
                            position={[currentData?.event_address?.coordinates?.lat ?? 0, currentData?.event_address?.coordinates?.lon ?? 0]}
                            selectedLocation={{
                                display_name: currentData?.event_address?.address || "",
                                lat: String(currentData?.event_address?.coordinates?.lat),
                                lon: String(currentData?.event_address?.coordinates?.lon)
                            }} />}
                        proposedContent={<MapPicker
                            position={[updateData?.updated_address?.coordinates?.lat ?? 0, updateData?.updated_address?.coordinates?.lon ?? 0]}
                            selectedLocation={{
                                display_name: updateData?.updated_address?.address || "",
                                lat: String(updateData?.updated_address?.coordinates?.lat),
                                lon: String(updateData?.updated_address?.coordinates?.lon)
                            }} />}
                    />
                </DetailCard>

                <Modal
                    isOpen={openBanner}
                    onClose={() => setOpenBanner(false)}
                >
                    <ViewBannerModal
                        bannerUrl={currentData?.event?.banner || ""}
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

    );
}

export function DateTimeBlock({ date, timeRange }: { date: string, timeRange: string }) {
    return (
        <div className="space-y-1">
            <div className="flex items-center gap-1.5">
                <CalenderIcon className="w-3.5 h-3.5 text-gray-400 shrink-0" />
                <span>{date}</span>
            </div>
            <div className="flex items-center gap-1.5">
                <ClockIcon className="w-3.5 h-3.5 text-gray-400 shrink-0" />
                <span>
                    {timeRange}
                    <span className="text-gray-400 ml-1">WIB</span>

                </span>
            </div>
        </div>
    );
}

function DescriptionBlock({ description }: {
    description: string | undefined;
}) {
    return (
        <div
            className="prose prose-sm text-sm text-gray-600 leading-relaxed space-y-3"
            dangerouslySetInnerHTML={{ __html: description || "" }}
        />
    );
}

function LocationBlock({
    loc,
}: {
    loc: {
        address: string | undefined;
        city: string | undefined;
        country: string | undefined;
        detailAddress: string | undefined;
    };
}) {
    return (
        <div className="space-y-2">
            <div className="grid grid-cols-2 gap-x-4 gap-y-1">
                <div>
                    <span className="text-[10px] uppercase tracking-widest text-gray-400 block">
                        Address
                    </span>
                    <span className="text-sm">{loc.address}</span>
                </div>

                <div>
                    <span className="text-[10px] uppercase tracking-widest text-gray-400 block">
                        City, Country
                    </span>
                    <span className="text-sm">
                        {loc.city}, {loc.country}
                    </span>
                </div>
            </div>

            <div>
                <span className="text-[10px] uppercase tracking-widest text-gray-400 block">
                    Detail Address
                </span>
                <span className="text-sm">{loc.detailAddress}</span>
            </div>
        </div>
    );
}

