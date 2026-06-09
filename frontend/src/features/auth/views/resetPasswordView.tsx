import { IoIosArrowRoundBack } from "react-icons/io";
import AuthContainer from "../components/container";
import Link from "../../../components/link";
import ResetPasswordForm from "../components/resetPasswordForm";
import { useResetPassword } from "../hooks/useResetPassword";

export default function ResetPasswordView() {
    const baseUrlPort = import.meta.env.VITE_URL_PORT;

    const {
        resetPassword,
        loading,
        message,
        error,
        errors,
    } = useResetPassword();

    const handleResetPassword = async (
        new_password: string,
        confirm_password: string
    ) => {
        await resetPassword({
            new_password,
            confirm_password,
        });
    };

    return (
        <AuthContainer className="max-w-lg w-full">
            <div className="w-full px-10 py-10">
                <div className="flex items-center gap-2">
                    <IoIosArrowRoundBack className="text-blue-500 text-xl" />

                    <Link endpoint={`http://localhost:${baseUrlPort}/login`}>
                        Back to Login
                    </Link>
                </div>

                <h2 className="mt-6 font-bold text-3xl text-center">
                    Reset Password
                </h2>

                <p className="text-center text-gray-500 mt-2 mb-6">
                    Create a new password for your account.
                </p>

                <ResetPasswordForm
                    onSubmit={handleResetPassword}
                    loading={loading}
                    errors={errors}
                />

                {message && (
                    <p className="text-green-600 mt-4 text-center">
                        {message}
                    </p>
                )}

                {error && (
                    <p className="text-red-500 mt-4 text-center">
                        {error}
                    </p>
                )}
            </div>
        </AuthContainer>
    );
}