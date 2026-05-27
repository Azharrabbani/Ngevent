import { useNavigate } from "react-router-dom";
import CompleteProfileContainer from "../components/container";
import { useCreateAttendeeProfile } from "../hooks/attendee/useCreateAttendeeProfile";
import type { CreateAttendeeProfileReq, CreateOrganizerProfileReq } from "../types/profileRequest";
import { useAuth } from "../../../lib/auth";
import { useCreateOrganizerProfile } from "../hooks/organizer/useCreateOrganizerProfile";
import CompleteAttendeeProfileForm from "../components/completeAttendeeProfileForm";
import CompleteOrganizerProfileForm from "../components/completeOrganizerProfileForm";

export default function CompleteProfileView() {
    const { mutateAsync: attendee, isPending: attendeePending, error: attendeeErrors } = useCreateAttendeeProfile();

    const mapValidationErrors = (error: any) => {
        const validationError = error?.response?.data?.error;

        if (!Array.isArray(validationError)) return {};

        return validationError.reduce((acc: any, e: any) => {
            acc[e.field] = e.message;
            return acc;
        }, {});
    };

    const attendeeValidationErrors = mapValidationErrors(attendeeErrors);

    const {
        mutateAsync: organizer,
        isPending: organizerLoading,
        error: organizerError,
    } = useCreateOrganizerProfile()

    const organizerValidationErrors = mapValidationErrors(organizerError);

    const { user, loading: userLoading } = useAuth()

    const navigate = useNavigate();

    if (userLoading) {
        return <p className="text-center">Loading...</p>
    }

    const handleCreateAttendeeProfile = async (payload: CreateAttendeeProfileReq) => {
        await attendee(payload);
        navigate("/dashboard");
    }

    const handleCreateOrganizerProfile = async (payload: CreateOrganizerProfileReq) => {
        await organizer(payload);
        navigate("/profile");
    }

    return (
        <CompleteProfileContainer>
            <div className="w-full max-w-md md:max-w-lg mx-auto px-4 sm:px-6">

                {user?.role === "user" ? (
                    <>
                        <h1 className="mb-6 sm:mb-8 font-bold text-xl sm:text-2xl md:text-4xl text-center">
                            Let's complete your profile
                        </h1>
                        <CompleteAttendeeProfileForm
                            onSubmit={handleCreateAttendeeProfile}
                            loading={attendeePending}
                            errors={attendeeValidationErrors}
                        />
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

                        <CompleteOrganizerProfileForm
                            onSubmit={handleCreateOrganizerProfile}
                            loading={organizerLoading}
                            errors={organizerValidationErrors}
                        />

                        {organizerError && !Array.isArray((organizerError as any)?.response?.data?.error) && (
                            <p className="text-red-500 text-center mt-3">
                                {(organizerError as any)?.response?.data?.error || "Failed to create profile"}
                            </p>
                        )}
                    </>
                )}

            </div>
        </CompleteProfileContainer>
    )
}