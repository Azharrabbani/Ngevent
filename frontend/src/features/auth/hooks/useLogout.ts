import { useState } from "react"
import { logoutApi } from "../api/authApi"

export const useLogout = () => {
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState<string | null>(null)

    const logout = async () => {
        try {
            setLoading(true);
            setError(null);

            const res = await logoutApi();
            return res.data;
        } catch (err: any) {
            setError(err.response?.data?.error || "Logout failed")
            throw err; 
        } finally {
            setLoading(false);
        }
    }

    return { logout, loading, error };
}