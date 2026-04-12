import AuthContainer from "../components/container";
import GoogleButton from "../components/googleButton";
import LoginForm from "../components/loginForm";
import Link from "../../../components/link";
import { useLogin } from "../hooks/useLogin";
import { useNavigate } from "react-router-dom";

export default function LoginView() {
    const baseUrlPort = import.meta.env.VITE_URL_PORT

    const navigate = useNavigate()

    const {login, loading, error, errors} = useLogin()

    const handleLogin = async(email: string, password: string) => {
        const user = await login({email, password})

        if (!user) return

        if (user.role == null) {
            navigate("/select-role")
        } else {
            navigate("/dashboard")
        }

    }

    return(
        <AuthContainer>
            {/* Image */}
                <div className="md:block hidden w-1/2">
                    <img className="rounded-2xl h-132" 
                    src="https://plus.unsplash.com/premium_photo-1701760184917-38e25718ee3e?fm=jpg&q=60&w=3000&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8MXx8bG9naW4lMjBiYWNrZ3JvdW5kfGVufDB8fDB8fHww" 
                    alt="login-img" />

                </div>

                {/* Form */}
                <div className="md:w-1/2 px-8 md:px-16 sm:py-7">
                    <h2 className="font-bold text-2xl text-center">Login</h2>
                    <LoginForm onSubmit={handleLogin} loading={loading} errors={errors}/>
                    {error && <p className="text-red-500 mt-3">{error}</p>}
                    <div className="mt-6 grid grid-cols-3 items-center text-gray-400">
                        <hr  className="border-gray-400"/>
                        <p className="text-center text-sm">or</p>
                        <hr className="border-gray-400"/>
                    </div>

                    <GoogleButton/>
                    
                    <div className="mt-2 text-sm py-4 text-center">
                       <Link endpoint={`http://localhost:${baseUrlPort}/forget`}>
                            Forget password
                       </Link>
                    </div>

                    <hr className="border-gray-400"/>
                    <div className="my-3 text-sm px-3 flex justify-center items-center gap-2">
                        <p>If you don't have an account...</p>
                        <Link endpoint={`http://localhost:${baseUrlPort}/register`}>
                            Register
                       </Link>
                    </div>
                </div>
        </AuthContainer>
    )
}