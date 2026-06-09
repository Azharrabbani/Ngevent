import AuthContainer from "../components/container";
import GoogleButton from "../components/googleButton";
import RegisterForm from "../components/registerForm";
import Link from "../../../components/link";
import { useRegister } from "../hooks/useRegister";
import { useNavigate } from "react-router-dom";
import { IoIosArrowRoundBack } from "react-icons/io";

export default function RegisterView() {
    const baseUrlPort = import.meta.env.VITE_URL_PORT;
    const navigate = useNavigate();

    const {
        mutateAsync: register,
        isPending,
        error,
        errors,
        isSuccess,
    } = useRegister();

    const handleRegister = async (
        email: string,
        password: string,
        confirm_password: string,
    ) => {
        const user = await register({
            email,
            password,
            confirm_password,
        });

        if (user) {
            navigate("/verified-email");
        }
    };

    return (
        <AuthContainer className="max-w-lg w-full">
            <div className="w-full px-10 py-10">
                <div className="flex items-center gap-2 cursor-pointer">
                    <IoIosArrowRoundBack className="text-blue-500 text-xl" />

                    <Link endpoint={`http://localhost:${baseUrlPort}/`}>
                        Back to dashboard
                    </Link>
                </div>

                <h2 className="mt-6 font-bold text-3xl text-center">
                    Create Account
                </h2>

                <p className="text-center text-gray-500 mt-2 mb-6">
                    Sign up to become an Event Organizer on Ngevent
                </p>

                <RegisterForm
                    onSubmit={handleRegister}
                    loading={isPending}
                    errors={errors}
                />

                {isSuccess && (
                    <p className="text-green-500 mt-3 text-center">
                        Registration successful! Please check your email.
                    </p>
                )}

                {error && (
                    <p className="text-red-500 mt-3 text-center">
                        {(error as any)?.response?.data?.error}
                    </p>
                )}

                <div className="mt-6 grid grid-cols-3 items-center text-gray-400">
                    <hr className="border-gray-300" />
                    <p className="text-center text-sm">or</p>
                    <hr className="border-gray-300" />
                </div>

                <GoogleButton />

                <hr className="border-gray-200 my-6" />

                <div className="text-sm flex justify-center gap-2">
                    <p className="text-gray-500">
                        Already have an account?
                    </p>

                    <Link endpoint={`http://localhost:${baseUrlPort}/login`}>
                        Login
                    </Link>
                </div>
            </div>
        </AuthContainer>
    );
}