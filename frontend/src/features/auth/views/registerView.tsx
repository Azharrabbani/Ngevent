import AuthContainer from "../components/container";
import GoogleButton from "../components/googleButton";
import RegisterForm from "../components/registerForm";
import Link from "../../../components/link";
import { useRegister } from "../hooks/useRegister";

export default function RegisterView() {
    const baseUrlPort = import.meta.env.VITE_URL_PORT

    const { register, loading, message, error, errors } = useRegister()

    const handleRegister = async (email: string, password: string, confirm_password: string, role: string) => {
        const user = await register({ email, password, confirm_password, role })

        if (user) {
            window.open(`http://localhost:${baseUrlPort}/verified-email`)
        }
    }

    return (
        <AuthContainer>
            {/* Image */}
            <div className="md:block hidden w-1/2">
                <img className="rounded-2xl h-132"
                    src="https://plus.unsplash.com/premium_photo-1701760184917-38e25718ee3e?fm=jpg&q=60&w=3000&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MXx8bG9naW4lMjBiYWNrZ3JvdW5kfGVufDB8fDB8fHww"
                    alt="register-img" />

            </div>

            {/* Form */}
            <div className="md:w-1/2 my-6 px-8 md:px-16">
                <h2 className="font-bold text-2xl text-center">Register</h2>
                <RegisterForm onSubmit={handleRegister} loading={loading} errors={errors} />
                {message && <p>Registration successful! Please check your email for the verification code.</p>}
                {error && <p className="text-red-500 mt-3">{error}</p>}
                <div className="mt-6 grid grid-cols-3 items-center text-gray-400">
                    <hr className="border-gray-400" />
                    <p className="text-center text-sm">or</p>
                    <hr className="border-gray-400" />
                </div>
                <GoogleButton />

                <hr className="border-gray-400" />
                <div className="my-3 text-sm px-3 flex justify-center items-center gap-2">
                    <p>Already have an account...</p>
                    <Link endpoint={`http://localhost:${baseUrlPort}/login`}>
                        Login
                    </Link>
                </div>
            </div>
        </AuthContainer>
    )
}