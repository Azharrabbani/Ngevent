import { BrowserRouter, Routes, Route } from "react-router-dom"
import LoginView from "../views/auth/loginView"
export default function AppRoutes() {
    return(
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<LoginView />} />
            </Routes>
        </BrowserRouter>
    )
}