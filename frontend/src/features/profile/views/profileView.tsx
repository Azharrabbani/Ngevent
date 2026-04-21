import ProfileContainer from "../components/profileContainer"
import { useGetCurrentAttendeeProfile } from "../hooks/useGetCurrentAttendeeProfiel"
import AttendeeProfileForm from "../components/attendeeProfileForm";
import OrganizerProfileForm from "../components/organizerProfileForm";
import { useGetCurrentOrganizerProfile } from "../hooks/useGetCurrentOrganizerProfile";
import { useAuth } from "../../../lib/auth";

export default function ProfileView() {
    const { user, loading } = useAuth()

    const attendee = useGetCurrentAttendeeProfile(user?.role === "user")
    const organizer = useGetCurrentOrganizerProfile(user?.role === "event organizer")

    if (loading) return <div>Loading...</div>
    if (!user) return null

    if (user.role === "user") {
        if (attendee.isLoading) return <div>Loading...</div>

        return (
            <ProfileContainer role={user.role}>
                <h1 className="font-bold text-2xl mb-6">Profile</h1>
                <AttendeeProfileForm profile={attendee.data}/>
            </ProfileContainer>
        )
    }

    if (user.role === "event organizer") {
        if (organizer.isLoading) return <div>Loading...</div>

        return (
            <ProfileContainer role={user.role}>
                <h1 className={`font-bold text-2xl mb-6 flex items-center gap-1`}>
                    Profile
                    {organizer.data?.status.status === "approved" && 
                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="purple" className="bi bi-patch-check-fill" viewBox="0 0 16 16">
                        <path d="M10.067.87a2.89 2.89 0 0 0-4.134 0l-.622.638-.89-.011a2.89 2.89 0 0 0-2.924 2.924l.01.89-.636.622a2.89 2.89 0 0 0 0 4.134l.637.622-.011.89a2.89 2.89 0 0 0 2.924 2.924l.89-.01.622.636a2.89 2.89 0 0 0 4.134 0l.622-.637.89.011a2.89 2.89 0 0 0 2.924-2.924l-.01-.89.636-.622a2.89 2.89 0 0 0 0-4.134l-.637-.622.011-.89a2.89 2.89 0 0 0-2.924-2.924l-.89.01zm.287 5.984-3 3a.5.5 0 0 1-.708 0l-1.5-1.5a.5.5 0 1 1 .708-.708L7 8.793l2.646-2.647a.5.5 0 0 1 .708.708"/>
                        </svg>
                    }
                </h1>
                <OrganizerProfileForm profile={organizer?.data}/>
            </ProfileContainer>
        )
    }

    return null
}