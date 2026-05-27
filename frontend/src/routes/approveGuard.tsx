import { Navigate, Outlet } from "react-router-dom";
import { useGetCurrentOrganizerProfile } from "../features/profile/hooks/organizer/useGetCurrentOrganizerProfile";
import { useAuth } from "../lib/auth";

export default function ApproveGuard() {
    const { user, loading } = useAuth()
    const organizer = useGetCurrentOrganizerProfile(user?.role === "event organizer")

    if (loading || organizer.isLoading) return <div>Loading...</div>
    if (!user) return null;

    if (organizer.data?.status?.status !== "approved") {
        return <Navigate to="/profile" replace />
    }

    return <Outlet />
}