import { useState } from "react"
import type { ResendOtpRequest } from "../types/authRequest"
import { resendOtpApi } from "../api/authApi"
import toast from "react-hot-toast"
import { useNavigate } from "react-router-dom"

export const useResendOtp = () => {
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState<string | null>(null)
    const [errors, setErrors] = useState<Record<string, string>>({})
    const [message, setMessage] = useState<string | null>(null)
    const navigate = useNavigate();

    const resendOtp = async (payload: ResendOtpRequest) => {
        try {
            setLoading(true);
            setError(null);
            setErrors({});
            setMessage(null);

            const email = localStorage.getItem("verification_email");

            if (!email) {
                setError("Session expired, please register again");
                return null;
            }

            const res = await resendOtpApi({
                ...payload,
                email: String(email)
            });

            setMessage(res.data)
        } catch (err: any) {
            const validationError = err.response?.data?.error;

            if (Array.isArray(validationError)) {
                const formatedError: Record<string, string> = {};

                validationError.forEach((e: any) => {
                    formatedError[e.field] = e.message;
                })

                setErrors(formatedError);
            } else {
                toast.error(err.response?.data?.error || "Resend OTP failed");
                localStorage.removeItem("verification_email");
                navigate("/register");

            }
        } finally {
            setLoading(false);
        }
    }

    return { resendOtp, loading, message, error, errors };
}