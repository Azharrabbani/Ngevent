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
import RichTextEditor from "../../../components/richTextEditior";
import { useCloseOrganizerAccount } from "../hooks/organizer/useCloseOrganizerAccount";
import Modal from "../../../components/modal";


interface Props {
    profile: OrganizerResponse | undefined
};

export default function OrganizerProfileForm({ profile }: Props) {
    const [previewOpen, setPreviewOpen] = useState(false);
    const [previewNpwp, setPreviewNpwp] = useState(false);
    const [previewNib, setPreviewNib] = useState(false);
    const [npwpPreview, setNpwpPreview] = useState<string | undefined>();
    const [nibPreview, setNibPreview] = useState<string | undefined>();
    const [npwpBlobUrl, setNpwpBlobUrl] = useState<string | null>(null)
    const [nibBlobUrl, setNibBlobUrl] = useState<string | null>(null)
    const [showCloseModal, setShowCloseModal] = useState(false);

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
        watch,
        setValue,
        formState: { isDirty, errors }
    } = useForm<FormValues>({
        defaultValues: {
            description: ''
        }
    });

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

    const { mutateAsync: updatePhoto, isPending: isPhotoPending } = useUpdateOrganizerPhoto()
    const { mutateAsync: updateProfile, isPending: isProfilePending } = useUpdateOrganizerProfile();
    const { mutate: closeAccount, isPending: isClosePending } = useCloseOrganizerAccount();

    const handleUpdatePhoto = async (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (!file) return;
        await updatePhoto({ photo: file });
    }

    const navigate = useNavigate();

    const goBackHome = () => {
        navigate("/organizer/dashboard");
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

        if (data.npwpFile) payload.npwpFile = data.npwpFile;
        if (data.nibFile) payload.nibFile = data.nibFile;

        await updateProfile(payload);
    };

    const handleOpenNpwp = async () => {
        setPreviewNpwp(true)
        if (npwpBlobUrl) return

        const url = npwpPreview || profile?.company_detail?.npwp_file
        if (!url) return

        if (url.startsWith("blob:")) {
            setNpwpBlobUrl(url)
            return
        }

        const res = await fetch(url, { credentials: "include" })
        const blob = await res.blob()
        setNpwpBlobUrl(URL.createObjectURL(blob))
    }

    const handleOpenNib = async () => {
        setPreviewNib(true)
        if (nibBlobUrl) return

        const url = nibPreview || profile?.company_detail?.nib_file
        if (!url) return

        if (url.startsWith("blob:")) {
            setNibBlobUrl(url)
            return
        }

        const res = await fetch(url, { credentials: "include" })
        const blob = await res.blob()
        setNibBlobUrl(URL.createObjectURL(blob))
    }

    const handleConfirmCloseAccount = () => {
        closeAccount(undefined, {
            onSettled: () => setShowCloseModal(false),
        });
    };

    const status = profile?.status?.status;
    const requestUpdate = profile?.request_updates;
    const showCloseButton = status === "approved" && !requestUpdate;
    
    useEffect(() => {
        return () => {
            if (npwpBlobUrl && npwpBlobUrl.startsWith("blob:")) URL.revokeObjectURL(npwpBlobUrl)
        }
    }, [npwpBlobUrl])

    useEffect(() => {
        return () => {
            if (nibBlobUrl && nibBlobUrl.startsWith("blob:")) URL.revokeObjectURL(nibBlobUrl)
        }
    }, [nibBlobUrl])

    return (
        <>

            <Modal
                isOpen={showCloseModal}
                onClose={() => !isClosePending && setShowCloseModal(false)}
            >
                <div
                    className="relative bg-white rounded-2xl p-6 w-[90%] max-w-md shadow-xl flex flex-col gap-4"
                    onClick={(e) => e.stopPropagation()}
                >
                    <div className="flex flex-col gap-1">
                        <h2 className="text-lg font-bold text-gray-900">Close Account</h2>
                        <p className="text-sm text-gray-500">
                            Are you sure you want to close your account? This action is
                            permanent and cannot be undone. All your profile data, events,
                            and related records will be deactivated.
                        </p>
                    </div>

                    <div className="flex justify-end gap-3 mt-2">
                        <Button
                            type="button"
                            className="rounded-md px-4 text-gray-600 bg-white border border-gray-300 hover:bg-gray-50"
                            onClick={() => setShowCloseModal(false)}
                            disabled={isClosePending}
                        >
                            Cancel
                        </Button>
                        <Button
                            type="button"
                            className="rounded-md px-4 bg-red-600 hover:bg-red-700 text-white"
                            onClick={handleConfirmCloseAccount}
                            disabled={isClosePending}
                        >
                            {isClosePending ? "Closing..." : "Yes, Close My Account"}
                        </Button>
                    </div>
                </div>
            </Modal>


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
                                <path d="M15.502 1.94a.5.5 0 0 1 0 .706L14.459 3.69l-2-2L13.502.646a.5.5 0 0 1 .707 0l1.293 1.293z" />
                                <path d="M4.939 9.21l6.813-6.814 2 2-6.813 6.814a.5.5 0 0 1-.196.12l-2.414.805a.25.25 0 0 1-.316-.316l.805-2.414a.5.5 0 0 1 .121-.196z" />
                            </svg>
                        </div>
                    )}
                </UploadPhoto>

                {requestUpdate && (
                    <div className="space-y-0 text-center">
                        <p className="italic text-sm text-yellow-600">
                            Your update is currently being reviewed by the admin.
                        </p>
                        <p className="italic text-sm text-gray-500">
                            Please wait for the approval process to be completed.
                        </p>
                    </div>
                )}

                {status === "pending" && (
                    <div className="space-y-0 text-center">
                        <p className="italic text-sm text-yellow-600">
                            Your profile is currently being reviewed by the admin.
                        </p>
                        <p className="italic text-sm text-gray-500">
                            Please wait for the approval process to be completed.
                        </p>
                    </div>
                )}

                {status === "rejected" && (
                    <div className="space-y-0 text-center">
                        <p className="italic text-sm text-red-600">
                            Your profile has been rejected by the admin.
                        </p>
                        <p className="italic text-sm text-gray-500">
                            Please review the rejection reason on your email and update your profile information.
                        </p>
                    </div>
                )}

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
                    {...register("name", { required: "Name is required" })}
                    error={errors.name?.message}
                    placeholder="Enter your full name"
                />

                <div>
                    <label className="text-xs sm:text-sm font-medium mb-1 block" htmlFor="">
                        Description
                    </label>
                    <RichTextEditor
                        value={watch("description") ?? ""}
                        onChange={(val) => setValue("description", val, { shouldDirty: true })}
                        error={errors.description?.message}
                    />
                </div>

                <div>
                    <label className="text-xs sm:text-sm font-medium mb-1 block" htmlFor="">
                        Address
                    </label>
                    <textarea
                        rows={3}
                        className="w-full p-2 rounded-xl bg-gray-200 outline-none resize-none"
                        {...register("address")}
                        placeholder="Enter your residential address"
                    />
                </div>

                <div className="grid grid-cols-2 gap-6">
                    <Input
                        className="p-2 sm:p-3 text-sm sm:text-base"
                        label="Phone number"
                        type="text"
                        {...register("phonenumber", { required: "Phone number is required" })}
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
                        onlyNumber
                        {...register("npwp", { required: "NPWP number is required" })}
                        placeholder="00-000-000-00"
                        error={errors.npwp?.message}
                    />
                    <Input
                        className="p-2 sm:p-3 text-sm sm:text-base"
                        label="NIB number"
                        type="text"
                        onlyNumber
                        {...register("nib", { required: "NIB number is required" })}
                        placeholder="00-000-000-00"
                        error={errors.nib?.message}
                    />
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    {/* NPWP */}
                    <UploadFile
                        uniqueId="npwp-upload"
                        file={profile?.company_detail?.npwp_file}
                        onClickFile={handleOpenNpwp}
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
                                {npwpBlobUrl
                                    ? <iframe src={npwpBlobUrl} className="w-full h-full" />
                                    : <div className="flex items-center justify-center h-full text-gray-500">Loading...</div>
                                }
                            </div>
                        </div>
                    )}

                    {/* NIB */}
                    <UploadFile
                        uniqueId="nib-upload"
                        file={profile?.company_detail?.nib_file}
                        onClickFile={handleOpenNib}
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
                                {nibBlobUrl
                                    ? <iframe src={nibBlobUrl} className="w-full h-full" />
                                    : <div className="flex items-center justify-center h-full text-gray-500">Loading...</div>
                                }
                            </div>
                        </div>
                    )}
                </div>


                <div
                    className={`flex flex-col sm:flex-row gap-6 sm:gap-6 mt-2 ${showCloseButton ? "justify-between" : "justify-end"
                        }`}
                >
                    {showCloseButton && (
                        <Button
                            type="button"
                            className="rounded-md px-4 text-red-600 bg-white border border-red-300 hover:bg-red-50"
                            onClick={() => setShowCloseModal(true)}
                            disabled={isClosePending}
                        >
                            Close Account
                        </Button>
                    )}

                    <div className="flex flex-col sm:flex-row gap-3">
                        <Button
                            type="button"
                            className="rounded-md px-4 text-purple-500 bg-white hover:bg-[#FAFAFA]"
                            onClick={goBackHome}
                        >
                            Back To Home
                        </Button>

                        <Button
                            type="submit"
                            disabled={!isDirty || isProfilePending}
                            className="bg-[#312E81] hover:bg-[#432E81] rounded-md px-4"
                        >
                            {isProfilePending ? "Updating..." : "Save Changes"}
                        </Button>
                    </div>
                </div>
            </form>
        </>
    )
}