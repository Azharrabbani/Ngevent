import { useNavigate } from "react-router-dom";
import Button from "../../../components/Button";
import Input from "../../../components/input";
import UploadPhoto from "../../../components/uploadPhoto";
import type { OrganizerResponse } from "../types/profileResponse";
import { useEffect, useState } from "react";
import UploadFile from "../../../components/uploadFile";
import { useForm } from "react-hook-form";
import { useUpdateOrganizerProfile } from "../hooks/organizer/useUpdateOrganizerProfile";
import { GetIsoFromPhoneNumber } from "../utils/phoneNumber";
import toast from "react-hot-toast";
import { useUpdateOrganizerPhoto } from "../hooks/organizer/useUpdateOrganizerPhoto";


interface Props {
    profile: OrganizerResponse | undefined
 };

export default function OrganizerProfileForm({profile}: Props) {
    const [previewOpen, setPreviewOpen] = useState(false);
    const [previewNpwp, setPreviewNpwp] = useState(false);
    const [previewNib, setPreviewNib] = useState(false);
    const [npwpPreview, setNpwpPreview] = useState<string | undefined>();
    const [nibPreview, setNibPreview] = useState<string | undefined>();

    type FormValues = {
        name: string;
        description: string;
        address: string;
        phonenumber: string;
        email: string;
        instagram: string;
        npwp: string;
        npwpFile?: File;  
        nib: string;
        nibFile?: File;  
    };

    const {
        register,
        handleSubmit,
        reset,
        setValue,
        formState: {isDirty, errors}
    } = useForm<FormValues>();

    useEffect(() => {
        if (profile) {
            reset({
                name: profile.name,
                description: profile.company_detail.description,
                address: profile.address,
                phonenumber: profile.phone_number,
                email: profile.social_media.email,
                instagram: profile.social_media.instagram,
                npwp: profile.company_detail.npwp,
                nib: profile.company_detail.nib,
            });
        }
    }, [profile, reset])

    const {mutateAsync: updatePhoto, isPending: isPhotoPending} =  useUpdateOrganizerPhoto()
    const {mutateAsync: updateProfile, isPending: isProfilePending} = useUpdateOrganizerProfile();

    const handleUpdatePhoto = async (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];

        if (!file) return;

        await updatePhoto({photo: file});
    }

    const navigate = useNavigate();

    const goBackHome = () => {
        navigate("/dashboard");
    };

    const handleNpwpFile = (file: File) => {
        setValue("npwpFile", file, { shouldDirty: true });

        const previewUrl = URL.createObjectURL(file);
        setNpwpPreview(previewUrl);
    };

    const handleNibFile = (file: File) => {
        setValue("nibFile", file, { shouldDirty: true });

        const previewUrl = URL.createObjectURL(file);
        setNibPreview(previewUrl);
    };

    const handleUpdateProfile = async (data: FormValues) => {
        const iso = GetIsoFromPhoneNumber(data.phonenumber);
        if (!iso) {
            toast.error("Invalid phone number");
            return;
        }

        const payload: any = {
            name: data.name,
            description: data.description,
            phonenumber: data.phonenumber,
            address: data.address,
            email: data.email,
            instagram: data.instagram,
            npwp: data.npwp,
            nib: data.nib,
            iso,
        };

        if (data.npwpFile) {
            payload.npwpFile = data.npwpFile;
        }

        if (data.nibFile) {
            payload.nibFile = data.nibFile;
        }

        await updateProfile(payload);
    };

    return(
        <form
            onSubmit={handleSubmit(handleUpdateProfile)} 
            className="flex flex-col gap-4 sm:gap-5 md:gap-6">
                <UploadPhoto
                    className="mt-3 relative z-10 w-20 h-20 sm:w-24 sm:h-24 md:w-28 md:h-28 rounded-full"
                    onClickImage={() => setPreviewOpen(true)}
                    onChange={handleUpdatePhoto}
                    disabled={isPhotoPending}
                    showEditIcon
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
                {profile?.status.status === "pending" ? 
                    <p className="italic text-sm text-center">your profile is being review by the admin</p> :
                    profile?.status.status === "rejected" ? 
                    <p className="italic text-sm text-center">your profile is being review by the admin</p> : ""
                }

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
                
                <Input
                    className="p-2 sm:p-3 text-sm sm:text-base"
                    label="Name"
                    type="text"
                    {...register(
                        "name",
                        {required: "Name is required"} 
                    )}
                    error={errors.name?.message}
                    placeholder="Enter your full name"
                />
        
                <div>
                    <label 
                    className="text-xs sm:text-sm font-medium mb-1 block"
                    htmlFor="">
                        Description
                    </label>
                    <textarea
                    rows={3}
                    className="w-full p-2 sm:p-3 text-sm sm:text-base rounded-xl bg-gray-200 outline-none resize-none"                         
                    {...register("description")}
                    placeholder="Enter your residential address"
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
                    className={`w-full p-2 rounded-xl bg-gray-200 outline-none resize-none`}                         
                    {...register("address")}
                    placeholder="Enter your residential address"
                    />
                </div>
        
                <div className="grid grid-cols-2 gap-6">
                    <Input
                        className="p-2 sm:p-3 text-sm sm:text-base"
                        label="Phone number"
                        type="text"
                        {...register(
                            "phonenumber",
                            {required: "Phone number is required"}
                        )}
                        error={errors.phonenumber?.message}
                        placeholder="+1 (555) 000-0000"
                    />
                        
                    <Input
                        className="p-2 sm:p-3 text-sm sm:text-base"
                        label="Country"
                        type="text"
                        name="country"
                        value={profile?.country}
                        disabled
                    />
                </div>

                <div className="grid grid-cols-2 gap-6">
                    <Input
                        className="p-2 sm:p-3 text-sm sm:text-base"
                        label="Instagram"
                        type="text"
                        {...register("instagram", {
                            validate: (value) => {
                              if (!value) return true;
                            
                              const isValidUrl = /^https?:\/\/(www\.)?instagram\.com\/.+/.test(value);
                              return isValidUrl || "Must be a valid Instagram link";
                            },
                        })}
                        error={errors.instagram?.message}
                        placeholder="https://instagram.com/username"
                    />
                        
                    <Input
                        className="p-2 sm:p-3 text-sm sm:text-base"
                        label="Email"
                        type="text"
                        {...register("email")}
                        placeholder="Email"
                    />
                </div>

                <div className="grid grid-cols-2 gap-6">
                    <Input
                        className="p-2 sm:p-3 text-sm sm:text-base"
                        label="NPWP number"
                        type="text"
                        {...register(
                            "npwp",
                            {required: "NPWP number is required"}
                        )}
                        placeholder="00-000-000-00"
                    />
            
                    <Input
                        className="p-2 sm:p-3 text-sm sm:text-base"
                        label="NIB number"
                        type="text"
                        {...register(
                            "nib",
                            {required: "NIB number is required"}
                        )}
                        placeholder="00-000-000-00"
                    />
                </div>
        
        
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    {/* NPWP */}
                   <UploadFile
                    uniqueId="npwp-upload"
                    file={profile?.company_detail?.npwp_file}
                    onClickFile={() => setPreviewNpwp(true)}
                    onChange={handleNpwpFile}
                   >
                        NPWP File
                   </UploadFile>
        
                    {previewNpwp && (
                        <div
                            className="fixed inset-0 bg-black/70 flex items-center justify-center z-50"
                            onClick={() => setPreviewNpwp(false)}
                        >
                            <div
                                className="w-[90%] h-[90%] bg-white rounded-lg overflow-hidden"
                                onClick={(e) => e.stopPropagation()}
                            >
                                <iframe
                                    src={npwpPreview || profile?.company_detail?.npwp_file}
                                    className="w-full h-full"
                                />
                            </div>
                        </div>
                    )}
                    
                    {/* NIB */}
                    <UploadFile
                    uniqueId="nib-upload"
                    file={profile?.company_detail?.nib_file}
                    onClickFile={() => setPreviewNib(true)}
                    onChange={handleNibFile}
                   >
                        NIB File
                   </UploadFile>

                    {previewNib && (
                        <div
                            className="fixed inset-0 bg-black/70 flex items-center justify-center z-50"
                            onClick={() => setPreviewNib(false)}
                        >
                            <div
                                className="w-[90%] h-[90%] bg-white rounded-lg overflow-hidden"
                                onClick={(e) => e.stopPropagation()}
                            >
                                <iframe
                                    src={nibPreview || profile?.company_detail?.nib_file}
                                    className="w-full h-full"
                                />
                            </div>
                        </div>
                    )}
                </div>
                
                <div className="flex flex-col sm:flex-row 
                                gap-6 sm:gap-6 justify-end">
                    <Button
                    type="button" 
                    className="rounded-md px-4 text-purple-500 bg-white hover:bg-[#FAFAFA] "
                    onClick={goBackHome}
                    >
                        Back To Home
                    </Button>
            
                    <Button 
                    type="submit"
                    disabled={!isDirty || isProfilePending}
                    className="bg-[#312E81] hover:bg-[#432E81] rounded-md px-4">
                        {isProfilePending ? "Updating..." : "Save Changes"}   
                    </Button>
                </div>
        </form>
    )
}