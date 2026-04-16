import { useState } from "react"
import type { CreateOrganizerProfileReq } from "../types/profileRequest"
import { CreateOrganizerProfileApi } from "../api/profileApi"

export const useCreateOrganizerProfile = () => {
    const [loading, setLoading] = useState(false)
    const [message, setMessage] = useState<string | null>(null)
    const [error, setError] = useState<string | null>(null)
    const [errors, setErrors] = useState<Record<string, string>>({})

    const createProfile = async(payload: CreateOrganizerProfileReq) => {
        try {
            setLoading(true);
            setError(null);
            setErrors({});
            setMessage(null);

            const res = await CreateOrganizerProfileApi(payload);

            setMessage(res.data);
        } catch(err: any) {
            const validationError = err.response?.data?.error;

            if (Array.isArray(validationError)) {
                const formatedError: Record<string, string> = {};

                validationError.forEach((e: any) => {
                    formatedError[e.field] = e.message;
                })
                setErrors(formatedError);
                console.log("errors: ", formatedError);
            } else {
                setError(err.response?.data?.error || "Failed create the profile")
            }
        } finally {
            setLoading(false);
        }
    }

    return {createProfile, loading, message, error, errors};
}