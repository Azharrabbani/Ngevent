import { BrowserRouter, Routes, Route } from "react-router-dom"
import LoginView from "../features/auth/views/loginView"
import RegisterView from "../features/auth/views/registerView"
import ForgetPassword from "../features/auth/views/forgetPasswordView"
import ResetPasswordView from "../features/auth/views/resetPasswordView"
import VerifiedEmailView from "../features/auth/views/verifiedEmailView"
export default function AppRoutes() {
    return(
        <BrowserRouter>
            <Routes>
                <Route path="/login" element={<LoginView />} />
                <Route path="/register" element={<RegisterView/>} />
                <Route path="/forget" element={<ForgetPassword/>}/>
                <Route path="/reset-password" element={<ResetPasswordView/>}/>
                <Route path="/verified-email" element={<VerifiedEmailView/>}/>
            </Routes>
        </BrowserRouter>
    )
}