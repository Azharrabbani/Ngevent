import { useEffect, useState } from "react";
import Button from "../../../components/Button";
import Input from "../../../components/input";
import UploadPhoto from "../../../components/uploadPhoto";
import type { AttendeeResponse } from "../types/profileResponse";
import { useNavigate } from "react-router-dom";
import { useUpdateAttendeePhoto } from "../hooks/useUpdateAttendeePhoto";
import { useForm } from "react-hook-form";
import { GetIsoFromPhoneNumber } from "../utils/phoneNumber";
import { useUpdateAttendeeProfile } from "../hooks/useUpdateAttendeeProfile";


interface Props {
    profile: AttendeeResponse | undefined
};

export default function AttendeeProfileForm({profile}: Props) {
    type FormValues = {
        name: string;
        username: string;
        phonenumber: string;
        address: string;
    };

    const {
        register,
        handleSubmit,
        reset,
        formState: {isDirty, errors}
    } = useForm<FormValues>();

    useEffect(() => {
        if (profile) {
            reset({
                name: profile.name,
                username: profile.username,
                phonenumber: profile.phone_number,
                address: profile.address,
            });
        }
    }, [profile, reset]);

    const [previewOpen, setPreviewOpen] = useState(false);

    const { mutateAsync: updatePhoto, isPending: isPhotoPending} = useUpdateAttendeePhoto();
    const { mutateAsync: updateProfile, isPending: isProfilePending } = useUpdateAttendeeProfile();

    const handleUpdatePhoto = async(e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (!file) return;

        await updatePhoto({photo: file});
    };

    const handleUpdateProfile = async (data: FormValues) => {
        const iso = GetIsoFromPhoneNumber(data.phonenumber);
        if (!iso) {
            console.error("Invalid phone number, ISO not found");
            return;
        }

        const payload = {
            name: data.name,
            username: data.username,
            phone_number: data.phonenumber,
            address: data.address,
            iso,
        };

        await updateProfile(payload);
    };

    const navigate = useNavigate();

    const goBackHome = () => {
        navigate("/dashboard");
    };
    
    return(
        <form
            onSubmit={handleSubmit(handleUpdateProfile)} 
            className="flex flex-col gap-4 sm:gap-5 md:gap-6"
        >
            <UploadPhoto
                className="relative w-20 h-20 sm:w-24 sm:h-24 md:w-28 md:h-28 rounded-full"
                onClickImage={() => setPreviewOpen(true)}
                showEditIcon
                onChange={handleUpdatePhoto}
                disabled={isPhotoPending}
            >
                {profile?.photo_profile ? (
                    <img
                        src={profile.photo_profile}
                        className="w-full h-full object-cover rounded-full"
                    />
                ) : (
                    <div className="w-full h-full flex items-center justify-center">
                        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" viewBox="0 0 16 16">
                            <path d="M15.502 1.94a.5.5 0 0 1 0 .706L14.459 3.69l-2-2L13.502.646a.5.5 0 0 1 .707 0l1.293 1.293z"/>
                            <path d="M4.939 9.21l6.813-6.814 2 2-6.813 6.814a.5.5 0 0 1-.196.12l-2.414.805a.25.25 0 0 1-.316-.316l.805-2.414a.5.5 0 0 1 .121-.196z"/>
                        </svg>
                    </div>
                )}
            </UploadPhoto>

            {previewOpen && (
                <div
                    className="fixed inset-0 bg-black/70 flex items-center justify-center z-50"
                    onClick={() => setPreviewOpen(false)}
                >
                    <img
                        src={profile?.photo_profile}
                        className="max-w-[90%] max-h-[90%] rounded-lg"
                    />
                </div>
            )}

            <div className="absolute bottom-0 right-0 bg-white rounded-full p-1 shadow z-10">
                
            </div>
                            
            <Input
                className="p-2 sm:p-3 text-sm sm:text-base"
                label="Name"
                type="text"
                {...register(
                    "name",
                    {required: "name is required"},
                )}
                error={errors.name?.message}
                placeholder="Enter your full name"
            />
        
            <Input
                className="p-2 sm:p-3 text-sm sm:text-base"
                label="Email"
                type="text"
                name="email"
                disabled
                value={profile?.email}
            />
        
            <Input
                className="p-2 sm:p-3 text-sm sm:text-base"
                label="Username"
                type="text"
                {...register("username")}
                placeholder="Username"
            />
                            
            <div className="grid grid-cols-2 gap-6">
                <Input
                    className="p-2 sm:p-3 text-sm sm:text-base"
                    label="Phone number"
                    type="text"
                    {...register(
                        "phonenumber",
                        {
                            required: "Phone number is required",
                            validate: (value) => {
                                return GetIsoFromPhoneNumber(value) ? true : "Invalid phone number";
                            },
                        },
                    )}
                    error={errors.phonenumber?.message}
                    placeholder="+1 (555) 000-0000"
                />
        
                <Input
                    className="p-2 sm:p-3 text-sm sm:text-base"
                    label="Country"
                    type="text"
                    name="country"
                    disabled
                    value={profile?.country}
                />
            </div>
        
            <div>
                <label 
                className="text-xs sm:text-sm font-medium mb-1 block"
                htmlFor="">
                    Address
                </label>
                <textarea
                rows={3}
                className="w-full p-2 sm:p-3 text-sm sm:text-base rounded-xl bg-gray-200 outline-none resize-none"                         
                {...register("address")}
                placeholder="Enter your residential address"
                />
            </div>
        
            <div className="flex flex-col sm:flex-row 
                            gap-6 sm:gap-6 justify-end">
                <Button
                className="rounded-md px-4 text-[#0040A1] bg-white hover:bg-[#FAFAFA] "
                onClick={goBackHome}
                >
                    Back To Home
                </Button>
        
                <Button
                type="submit" 
                className="rounded-md px-4"
                disabled={!isDirty || isProfilePending}
                >
                    Save Changes
                </Button>
            </div>
        </form>        
    )
}