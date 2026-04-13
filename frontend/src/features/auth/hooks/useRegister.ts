import { useState } from "react"
import type { RegisterRequest } from "../types/authRequest"
import type { RegisterResponse } from "../types/authResponse"
import { registerApi } from "../api/authApi"
import type { successResponse } from "../../../types/apiResponse"

export const useRegister = () => {
    const [loading, setLoading] = useState(false)
    const [message, setMessage] = useState<string | null>(null)
    const [error, setError] = useState<string | null>(null)
    const [errors, setErrors] = useState<Record<string, string>>({})

    const register = async(payload: RegisterRequest): Promise<successResponse<RegisterResponse> | null> => {
        try{
            setLoading(true)
            setError(null)
            setErrors({})
            setMessage(null)

            const res = await registerApi(payload)

            localStorage.setItem("verify-email", res.data.email)

            setMessage(res.message);
            return res
        }catch(err: any) {
            const validationError = err.response?.data?.error

            if (Array.isArray(validationError)) {
                const formatedError: Record<string, string> = {}

                validationError.forEach((e: any) => {
                    formatedError[e.field] = e.message
                })

                setErrors(formatedError)
            } else {
                setError(err.response?.data?.error || "Register failed")
            }

            return null
        }finally {
            setLoading(false)
        }
    }

    return {register, loading, message, error, errors}
}