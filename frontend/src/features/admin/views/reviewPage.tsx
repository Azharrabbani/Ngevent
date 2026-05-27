import { useParams, useNavigate, useSearchParams } from "react-router-dom";
import { useGetUpdateByEventID } from "../../events/hooks/useGetUpdateByEventID";
import { useReviewEvent } from "../../events/hooks/useReviewEvent";

import UpdateRequestView from "../components/events/updateRequest";
import EventDetailView from "../components/events/eventDetail";
import AdminSidebar from "../components/sideBar";
import { useReviewUpdatedEvent } from "../../events/hooks/useReviewUpdateEvent";

export default function ReviewPage() {
    const { id } = useParams<{ id: string }>();
    const [searchParams] = useSearchParams();
    const status = searchParams.get("status") || "";
    const navigate = useNavigate();

    const { data: updateData, isLoading: updateLoading, isError } = useGetUpdateByEventID(id!, status);
    const { mutateAsync: reviewEvent, isPending: reviewPending } = useReviewEvent();
    const { mutateAsync: reviewUpdated, isPending: reviewUpdatedPending } = useReviewUpdatedEvent();

    if (updateLoading) {
        return (
            <AdminSidebar>
                <div className="flex items-center justify-center h-screen">
                    <p className="text-gray-400 text-sm">Loading...</p>
                </div>
            </AdminSidebar>
        );
    }

    const hasUpdate = !isError && updateData != null;

    if (hasUpdate) {
        return (
            <AdminSidebar>
                <UpdateRequestView
                    status={status}
                    onBack={() => navigate(-1)}
                    isSubmitting={reviewUpdatedPending}
                    onApprove={async () => {
                        await reviewUpdated({ id: updateData.id, status: "approved" });
                        navigate(-1);
                    }}
                    onReject={async (reason: string) => {
                        await reviewUpdated({ id: updateData.id, status: "rejected", reason });
                        navigate(-1);
                    }}
                />
            </AdminSidebar>
        );
    }

    return (
        <AdminSidebar>
            <EventDetailView
                id={id}
                onBack={() => navigate(-1)}
                isSubmitting={reviewPending}
                onApprove={async () => {
                    await reviewEvent({ id: id!, status: "active" });
                    navigate(-1);
                }}
                onReject={async (reason: string) => {
                    await reviewEvent({ id: id!, status: "rejected", reason });
                    navigate(-1);
                }}
            />
        </AdminSidebar>
    );
}