import { useState } from "react"
import type { ResetPasswordRequest } from "../types/authRequest"
import { resetPasswordApi } from "../api/authApi"
import { useSearchParams } from "react-router-dom"

export const useResetPassword = () => {
    const [searchParams] = useSearchParams();

    const [loading, setLoading] = useState(false)
    const [error, setError] = useState<string | null>(null)
    const [errors, setErrors] = useState<Record<string, string>>({})
    const [message, setMessage] = useState<string | null>(null)

    const resetPassword = async(payload: ResetPasswordRequest) => {
        try {
            setLoading(true)
            setError(null)
            setErrors({})
            setMessage(null)
            
            const token = searchParams.get("token")
            
            const res = await resetPasswordApi(payload, String(token))
            
            setMessage(res.data)
        }catch(err: any) {
            const validationErrors = err.response?.data?.error

            if (Array.isArray(validationErrors)) {
                const formatedError: Record<string, string> = {}
                
                validationErrors.forEach((e: any) => {
                    formatedError[e.field] = e.message
                })

                setErrors(formatedError)
            } else {
                setError(err.response?.data?.error || "failed reset your password")
            }
        }finally {
            setLoading(false)
        }
    }

    return {resetPassword, loading, message, error, errors}
}