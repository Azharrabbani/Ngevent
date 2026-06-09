import { IoIosArrowRoundBack } from "react-icons/io";
import AuthContainer from "../components/container";
import Link from "../../../components/link";
import ForgetPasswordForm from "../components/forgetPasswordForm";
import { useForgetPassword } from "../hooks/useForgetPassword";

export default function ForgetPassword() {
    const baseUrlPort = import.meta.env.VITE_URL_PORT;

    const {
        forgetPassword,
        loading,
        message,
        error,
        errors,
    } = useForgetPassword();

    const handleForgetPassword = async (email: string) => {
        await forgetPassword({ email });
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
                    Forgot Password
                </h2>

                <p className="text-center text-gray-500 mt-2 mb-6">
                    Enter your email address and we'll send you a password reset link.
                </p>

                <ForgetPasswordForm
                    onSubmit={handleForgetPassword}
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