import { useNavigate } from "react-router-dom";
import AttendeeProfileForm from "../components/attendeeProfileForm";
import CompleteProfileContainer from "../components/container";
import { UseCreateAttendeeProfile } from "../hooks/useCreateAttendeeProfile";
import type { CreateAttendeeProfileReq, CreateOrganizerProfileReq } from "../types/profileRequest";
import { useAuth } from "../../../lib/auth";
import OrganizerProfileForm from "../components/organizerProfileForm";
import { useCreateOrganizerProfile } from "../hooks/useCreateOrganizerProfile";

export default function CompleteProfileView() {
    const { 
        createProfile: attendee, 
        loading: attendeeLoading, 
        message: attendeeMessage, 
        error: attendeeError, 
        errors: attendeeErrors } = UseCreateAttendeeProfile();
    
    const {
        createProfile: organizer,
        loading: organizerLoading,
        message: organizerMessage,
        error: organizerError,
        errors: organizerErrors
    } = useCreateOrganizerProfile();
    
    
    const { user, loading: userLoading } = useAuth()
    
    const navigate = useNavigate();

    if (userLoading) {
        return <p className="text-center">Loading...</p>
    }

    const handleCreateAttendeeProfile = async(payload: CreateAttendeeProfileReq) => {
        await attendee(payload);
        navigate("/dashboard");
    }

    const handleCreateOrganizerProfile = async(payload: CreateOrganizerProfileReq) => {
        await organizer(payload);
        navigate("/dashboard");
    }

    return(
        <CompleteProfileContainer>
            <div className="w-full max-w-md md:max-w-lg mx-auto px-4 sm:px-6">

                {user?.role === "user" ? (
                    <>
                        <h1 className="mb-6 sm:mb-8 font-bold text-xl sm:text-2xl md:text-4xl text-center">
                            Let's complete your profile
                        </h1>    
                        <AttendeeProfileForm 
                            onSubmit={handleCreateAttendeeProfile} 
                            loading={attendeeLoading} 
                            errors={attendeeErrors}
                        />

                        {attendeeError && (
                            <p className="text-red-500 text-center mt-3">{attendeeError}</p>
                        )}

                        {attendeeMessage && (
                            <p className="text-green-500 text-center mt-3">{attendeeMessage}</p>
                        )}
                    </>
                ) : (
                    <>
                        <h1 className="mb-6 sm:mb-8 font-bold text-xl sm:text-2xl md:text-4xl text-center">
                        Let's complete your profile
                        </h1>
                        <p>
                            This information helps us verify your business and
                            personalize your security experience.
                        </p>
                        
                        <OrganizerProfileForm
                        onSubmit={handleCreateOrganizerProfile} 
                        loading={organizerLoading} 
                        errors={organizerErrors}/>

                        {organizerError && (
                            <p className="text-red-500 text-center mt-3">{organizerError}</p>
                        )}

                        {organizerMessage && (
                            <p className="text-green-500 text-center mt-3">{organizerMessage}</p>
                        )}
                    </>
                )}
                
            </div>
        </CompleteProfileContainer>
    )
}