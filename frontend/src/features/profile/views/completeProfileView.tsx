import { useNavigate } from "react-router-dom";
import AttendeeProfileForm from "../components/attendeeProfileForm";
import CompleteProfileContainer from "../components/container";
import { UseCreateAttendeeProfile } from "../hooks/useCreateAttendeeProfile";
import type { CreateAttendeeProfileReq } from "../types/profileRequest";

export default function CompleteProfileView() {
    const {createProfile, loading, message, error, errors} = UseCreateAttendeeProfile();
    
    const navigate = useNavigate();

    const handleCreateProfile = async(payload: CreateAttendeeProfileReq) => {
        await createProfile(payload);
        navigate("/dashboard");
    }

    return(
        <CompleteProfileContainer>
            <div className="w-full max-w-md md:max-w-lg mx-auto px-4 sm:px-6">
                <h1 className="mb-6 sm:mb-8 font-bold text-xl sm:text-2xl md:text-4xl text-center">
                    Let's complete your profile
                </h1>    

                <AttendeeProfileForm onSubmit={handleCreateProfile} loading={loading} errors={errors}/>
                {error && (
                    <p className="text-red-500 text-center mt-3">{error}</p>
                )}

                {message && (
                    <p className="text-green-500 text-center mt-3">{message}</p>
                )}
            </div>
        </CompleteProfileContainer>
    )
}