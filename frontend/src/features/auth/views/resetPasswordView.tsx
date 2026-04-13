import AuthContainer from "../components/container";
import ResetPasswordForm from "../components/resetPasswordForm";
import { useResetPassword } from "../hooks/useResetPassword";

export default function ResetPasswordView() {
    const {resetPassword, loading, message, error, errors} = useResetPassword()

    const handleResetPassword = async(new_password: string, confirm_password: string) => {
        await resetPassword({new_password, confirm_password})
    }

    return (
        <AuthContainer>
            <div className="md:block hidden w-1/2">
                <img className="rounded-2xl h-132" 
                src="https://images-wixmp-ed30a86b8c4ca887773594c2.wixmp.com/i/2327ffa4-20e9-41b6-b77c-462d26f7bfea/d3ecbzg-045f95bb-8610-4dbb-84ab-e1d08fb037d7.jpg" 
                alt="reset-password-img" />
            </div>

            {/* Form */}
            <div className="md:w-1/2 px-8 md:px-16 sm:py-7">
                <h2 className="font-bold text-2xl text-center">Reset passsword</h2>
                <ResetPasswordForm onSubmit={handleResetPassword} loading={loading} errors={errors}/>
                {message && <p className="text-green-400 mt-3">{message}</p>}
                {error && <p className="text-red-500 mt-3">{error}</p>}
            </div>
        </AuthContainer>
    )
}