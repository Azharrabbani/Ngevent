import ProfileContainer from "../components/profileContainer"
import { useGetCurrentAttendeeProfile } from "../hooks/useGetCurrentAttendeeProfiel"
import AttendeeProfileForm from "../components/attendeeProfileForm";

export default function ProfileView() {
    const {profile, loading} = useGetCurrentAttendeeProfile();

    if (loading) return <div>Loading...</div>;
    if (!profile) return null;
    
    return (
        <ProfileContainer>
            <h1 className="font-bold text-2xl sm:text-3xl mb-6 sm:mb-8">
                Profile
            </h1>

            <AttendeeProfileForm profile={profile}/>

            {/* Next organizer profile */}
        </ProfileContainer>
    )
}