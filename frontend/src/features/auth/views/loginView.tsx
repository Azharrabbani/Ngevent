import { useNavigate } from "react-router-dom";
import { useLogin } from "../hooks/useLogin";
import AuthContainer from "../components/container";
import { IoIosArrowRoundBack } from "react-icons/io";
import Link from "../../../components/link";
import LoginForm from "../components/loginForm";
import GoogleButton from "../components/googleButton";

export default function LoginView() {
    const baseUrlPort = import.meta.env.VITE_URL_PORT;

    const navigate = useNavigate();

    const { login, loading, error, errors } = useLogin();

    const handleLogin = async (email: string, password: string) => {
        const user = await login({ email, password });

        if (!user) return;

        if (user.role === "admin") {
            navigate("/admin/dashboard");
        }

        if (user.role === "event organizer") {
            navigate("/organizer/dashboard");
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
                    Welcome Back!
                </h2>

                <p className="text-center text-gray-500 mt-2 mb-6">
                    Sign in to your Ngevent account
                </p>

                <LoginForm
                    onSubmit={handleLogin}
                    loading={loading}
                    errors={errors}
                />

                <div className="flex justify-end mt-2">
                    <Link
                        endpoint={`http://localhost:${baseUrlPort}/forget`}
                        className="text-sm text-blue-500 hover:text-blue-600"
                    >
                        Forgot Password?
                    </Link>
                </div>

                {error && (
                    <p className="text-red-500 mt-3 text-center">
                        {error}
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
                        Don't have an account?
                    </p>

                    <Link endpoint={`http://localhost:${baseUrlPort}/register`}>
                        Register
                    </Link>
                </div>
            </div>
        </AuthContainer>
    );
}