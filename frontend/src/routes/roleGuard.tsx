import { Navigate, Outlet } from "react-router-dom"
import { useAuth } from "../lib/auth"

export default function RoleGuard() {
    const { user, loading } = useAuth()

    if (loading) return <p>Loading...</p>

    if (!user?.role) {
        return <Navigate to="/select-role" replace />
    }

    return <Outlet />
}