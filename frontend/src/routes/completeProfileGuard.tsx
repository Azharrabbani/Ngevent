import { Navigate, Outlet } from "react-router-dom"
import { useAuth } from "../lib/auth"

export default function CompleteProfileGuard() {
    const { user, loading } = useAuth()

    if (loading) return <p>Loading...</p>

    if (user?.has_profile) {
        return <Navigate to="/profile" replace />
    }

    return <Outlet />
}