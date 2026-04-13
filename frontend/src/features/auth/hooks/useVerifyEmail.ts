import { useState } from "react"
import type { VerifyEmailRequest } from "../types/authRequest"
import { verifyEmailApi } from "../api/authApi"

export const useVerifyEmail = () => {
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState<string | null>(null)
    const [errors, setErrors] = useState<Record<string, string>>({})
    const [message, setMessage] = useState<string | null>(null)
    const baseUrlPort = import.meta.env.VITE_URL_PORT

    const verifyEmail = async(payload: VerifyEmailRequest) => {
        try {
            setLoading(true)
            setError(null)
            setErrors({})
            setMessage(null)

            const email = localStorage.getItem("verify-email")

            if (!email) {
                setError("Session expired, please register again")
                return null
            }

            const res = await verifyEmailApi({
                ...payload,
                email
            })

            setMessage(res.data)

            localStorage.removeItem("verify-email")

            window.location.href = `http://localhost:${baseUrlPort}/login`

        } catch(err: any) {
            const validationError = err.response?.data?.error

            if (Array.isArray(validationError)) {
                const formatedError: Record<string, string> = {}

                validationError.forEach((e: any) => {
                    formatedError[e.field] = e.message
                })

                setErrors(formatedError)
            } else {
                setError(err.response?.data?.error || "Verify email failed")
            }

            return null
        } finally {
            setLoading(false)
        }
    }

    return {verifyEmail, loading, message, error, errors}
}