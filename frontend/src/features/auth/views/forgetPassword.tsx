import AuthContainer from "../components/container";
import Link from "../../../components/link";
import ForgetPasswordForm from "../components/forgetPasswordForm";

export default function ForgetPassword() {
    const baseUrlPort = import.meta.env.VITE_URL_PORT

    return (
        <AuthContainer>
            <div className="md:block hidden w-1/2">
                <img className="rounded-2xl h-132" 
                src="https://images-wixmp-ed30a86b8c4ca887773594c2.wixmp.com/i/2327ffa4-20e9-41b6-b77c-462d26f7bfea/d3ecbzg-045f95bb-8610-4dbb-84ab-e1d08fb037d7.jpg" 
                alt="login-img" />
            </div>

            {/* Form */}
            <div className="md:w-1/2 px-8 md:px-16 sm:py-7">
                <div className="mb-8 flex gap-1 items-center">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" className="bi bi-arrow-left text-blue-500" viewBox="0 0 16 16">
                      <path fill-rule="evenodd" d="M15 8a.5.5 0 0 0-.5-.5H2.707l3.147-3.146a.5.5 0 1 0-.708-.708l-4 4a.5.5 0 0 0 0 .708l4 4a.5.5 0 0 0 .708-.708L2.707 8.5H14.5A.5.5 0 0 0 15 8"/>
                    </svg>
                    <Link endpoint={`http://localhost:${baseUrlPort}/login`}>
                        Back
                    </Link>
                </div>
                
                <h2 className="font-bold text-2xl text-center">Forget passsword</h2>
                <ForgetPasswordForm/>
            </div>
        </AuthContainer>
    )
}