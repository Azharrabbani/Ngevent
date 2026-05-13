import { BrowserRouter, Routes, Route } from "react-router-dom"
import LoginView from "../features/auth/views/loginView"
import RegisterView from "../features/auth/views/registerView"
import ForgetPassword from "../features/auth/views/forgetPasswordView"
import ResetPasswordView from "../features/auth/views/resetPasswordView"
import VerifiedEmailView from "../features/auth/views/verifiedEmailView"
import SelectRoleView from "../features/auth/views/selectRoleView"
import ProtectedRoute from "./protected"
import CompleteProfileView from "../features/profile/views/completeProfileView"
import CompleteProfileGuard from "./completeProfileGuard"
import RoleGuard from "./roleGuard"
import ProfileView from "../features/profile/views/profileView"
import ProfileGuard from "./profileGuard"
import AdminDashboard from "../features/admin/views/dashboard"
import AttendeeList from "../features/admin/views/attendeeList"
import OrganizerList from "../features/admin/views/organizerList"
import ActiveEventList from "../features/admin/views/activeEventList"
import PendingEventList from "../features/admin/views/pendingEventList"
import OrganizerDashboard from "../layout/organizerDashboard"
import CheckEvents from "../features/events/organizer/views/checkEvents"
import CancelEvent from "../features/events/organizer/views/cancelEvent"
import DraftEvents from "../features/events/organizer/views/draftEvent"
import EventFormView from "../features/events/organizer/views/eventFormView"
import ApproveGuard from "./approveGuard"

export default function AppRoutes() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/login" element={<LoginView />} />
                <Route path="/register" element={<RegisterView />} />
                <Route path="/forget" element={<ForgetPassword />} />
                <Route path="/reset-password" element={<ResetPasswordView />} />
                <Route path="/verified-email" element={<VerifiedEmailView />} />

                <Route element={<ProtectedRoute />}>
                    <Route path="/select-role" element={<SelectRoleView />} />

                    <Route element={<RoleGuard allowedRoles={["user", "event organizer"]} />}>

                        <Route element={<CompleteProfileGuard />}>
                            <Route path="/complete-profile" element={<CompleteProfileView />} />
                        </Route>

                        <Route element={<ProfileGuard />}>
                            <Route path="/profile" element={<ProfileView />} />
                        </Route>
                    </Route>
                </Route>

                <Route path="/admin" element={<ProtectedRoute />}>
                    <Route element={<RoleGuard allowedRoles={["admin"]} />}>
                        <Route path="/admin/dashboard" element={<AdminDashboard />} />
                        <Route path="/admin/organizer-list" element={<OrganizerList />} />
                        <Route path="/admin/attendee-list" element={<AttendeeList />} />
                        <Route path="/admin/events-active" element={<ActiveEventList />} />
                        <Route path="/admin/events-pending" element={<PendingEventList />} />
                    </Route>
                </Route>

                <Route path="/organizer" element={<ProtectedRoute />}>
                    <Route element={<RoleGuard allowedRoles={["event organizer"]} />}>
                        <Route path="dashboard" element={
                            <OrganizerDashboard>
                                <CheckEvents />
                            </OrganizerDashboard>
                        }
                        />

                        <Route path="cancel-event" element={
                            <OrganizerDashboard>
                                <CancelEvent />
                            </OrganizerDashboard>
                        }
                        />

                        <Route path="draft-event" element={
                            <OrganizerDashboard>
                                <DraftEvents />
                            </OrganizerDashboard>
                        }
                        />


                        <Route element={<ProfileGuard />}>
                            <Route element={<ApproveGuard />}>
                                <Route path="event/new" element={<EventFormView />} />
                                <Route path="event/edit/:id" element={<EventFormView />} />
                            </Route>
                        </Route>

                    </Route>
                </Route>
            </Routes>
        </BrowserRouter>
    )
}