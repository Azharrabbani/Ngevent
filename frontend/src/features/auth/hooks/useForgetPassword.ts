import { useState } from "react"
import type { ForgetPasswordRequest } from "../types/authRequest"
import { forgetPasswordApi } from "../api/authApi"

export const useForgetPassword = () => {
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState<string | null>(null)
    const [errors, setErrors] = useState<Record<string, string>>({})
    const [message, setMessage] = useState<string | null>(null)

    const forgetPassword = async(payload: ForgetPasswordRequest) => {
        try {
            setLoading(true)
            setError(null)

            const res = await forgetPasswordApi(payload)

            setMessage(res.data)
        } catch(err: any) {
            setError(null)
            setErrors({})
            setMessage(null)
            
            const validationErrors = err.response?.data?.error

            if (Array.isArray(validationErrors)) {
                const formatedError: Record<string, string> = {}

                validationErrors.forEach((e: any) => {
                    formatedError[e.field] = e.message
                })

                setErrors(formatedError)
            } else {
                setError(err.response?.data?.error || "failed send email")
            }
        } finally {
            setLoading(false)
        }
    }

    return {forgetPassword, loading, message, error, errors}
}