import { Navigate, Outlet } from "react-router-dom"
import { useAuth } from "../lib/auth"

type RoleGuardProps = {
    allowedRoles: string[]
}

export default function RoleGuard({ allowedRoles }: RoleGuardProps) {
    const { user, loading } = useAuth()

    if (loading) return <p>Loading...</p>

    if (!user?.role) {
        return <Navigate to="/select-role" replace />
    }

    if (!allowedRoles.includes(user.role)) {
        return <Navigate to="/unauthorized" replace />
    }

    return <Outlet />
}