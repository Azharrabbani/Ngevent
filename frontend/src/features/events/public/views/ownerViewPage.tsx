import { useParams } from "react-router-dom";
import { useGetOrganizerBySlug } from "../../../profile/hooks/organizer/useGetOrganizerBySlug"
import { useGetPublicOrganizerEvents } from "../../hooks/useGetPublicOrganizerEvents";

export default function OwnerViewPage() {
    const { slug } = useParams<{ slug: string }>();

    const { data: owner, isLoading: ownerLoading } =
        useGetOrganizerBySlug(slug!);

    const { data: events, isLoading: eventsLoading } =
        useGetPublicOrganizerEvents(owner?.id ?? "", {
            status: "active",
            pagination: {
                page: 1,
                limit: 8,
            },
        });

    if (ownerLoading) {
        return <div>Loading owner...</div>;
    }

    if (eventsLoading) {
        return <div>Loading events...</div>;
    }

    console.log(events);

    return (
        <h1>Event Owner</h1>
    );
}