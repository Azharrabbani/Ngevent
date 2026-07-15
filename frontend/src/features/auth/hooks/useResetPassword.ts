import { useState } from "react"
import type { ResetPasswordRequest } from "../types/authRequest"
import { resetPasswordApi } from "../api/authApi"
import { useNavigate, useSearchParams } from "react-router-dom"
import toast from "react-hot-toast"

export const useResetPassword = () => {
    const [searchParams] = useSearchParams();

    const [loading, setLoading] = useState(false)
    const [error, setError] = useState<string | null>(null)
    const [errors, setErrors] = useState<Record<string, string>>({})
    const [message, setMessage] = useState<string | null>(null)
    const navigate = useNavigate();

    const resetPassword = async (payload: ResetPasswordRequest) => {
        try {
            setLoading(true)
            setError(null)
            setErrors({})
            setMessage(null)

            const token = searchParams.get("token")

            const res = await resetPasswordApi(payload, String(token))

            setMessage(res.data)
        } catch (err: any) {
            const validationErrors = err.response?.data?.error

            if (Array.isArray(validationErrors)) {
                const formatedError: Record<string, string> = {}

                validationErrors.forEach((e: any) => {
                    formatedError[e.field] = e.message
                })

                setErrors(formatedError)
            } else {
                if (err.response?.data?.error === "Session has expired") {
                    toast.error(err.response?.data?.error);
                    navigate("/login");
                }
                setError(err.response?.data?.error || "failed reset your password")
            }
        } finally {
            setLoading(false)
        }
    }

    return { resetPassword, loading, message, error, errors }
}