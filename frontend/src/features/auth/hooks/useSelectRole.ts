import { useState } from "react"
import { selectRoleApi } from "../api/authApi"
import type { SelectRoleRequest } from "../types/authRequest"

export const useSelectRole = () => {
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [message, setMessage] = useState<string | null>(null);

    const selectRole = async(payload: SelectRoleRequest) => {
        try {
            setLoading(true);
            setError(null);
            setMessage(null)

            const res = await selectRoleApi(payload);

            setMessage(res.data);
        } catch(err: any) {
            setError(err.response?.data?.error || "failed to select the role")
        } finally {
            setLoading(false);
        }
    };

    return {selectRole, loading, error, message};
};