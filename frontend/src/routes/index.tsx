import { BrowserRouter, Routes, Route } from "react-router-dom"
import LoginView from "../views/auth/loginView"
import RegisterView from "../views/auth/registerView"
export default function AppRoutes() {
    return(
        <BrowserRouter>
            <Routes>
                <Route path="/login" element={<LoginView />} />
                <Route path="/register" element={<RegisterView/>} />
            </Routes>
        </BrowserRouter>
    )
}