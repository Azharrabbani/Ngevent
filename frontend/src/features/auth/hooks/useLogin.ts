import { useState } from "react"
import type { LoginRequest } from "../types/authRequest"
import type { AuthModel } from "../types/authModel"
import { loginApi } from "../api/authApi"
import { mapLoginResponse } from "../utils/mapResponse"

export const useLogin = () => {
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState<string | null>(null)
    const [errors, setErrors] = useState<Record<string, string>> ({})

    const login = async(payload: LoginRequest): Promise<AuthModel | null> => {
        try{
            setLoading(true)
            setError(null)
            setErrors({})

            const res = await loginApi(payload)

            const user = mapLoginResponse(res)

            localStorage.setItem("token", user.token)

            return user
        } catch(err: any) {
            const validationErrors = err.response?.data?.error

            if (Array.isArray(validationErrors)) {
                const formatedError: Record<string, string> = {}

                validationErrors.forEach((e: any) => {
                    formatedError[e.field] = e.message
                })

                setErrors(formatedError)
            } else {
                setError(err.response?.data?.error || "Login failed")
            }

            return null
        } finally {
            setLoading(false)
        }
    }

    return {login, loading, error, errors}
}