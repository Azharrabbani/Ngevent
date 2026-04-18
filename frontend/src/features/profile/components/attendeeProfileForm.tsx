import { useState } from "react";
import Button from "../../../components/Button";
import Input from "../../../components/input";
import UploadPhoto from "../../../components/uploadPhoto";
import type { AttendeeResponse } from "../types/profileResponse";

interface Props {
    profile: AttendeeResponse | null
};

export default function AttendeeProfileForm({profile}: Props) {
    const [previewOpen, setPreviewOpen] = useState(false);
    
    return(
        <form 
            className="flex flex-col gap-4 sm:gap-5 md:gap-6"
        >
            <UploadPhoto
                className="relative w-20 h-20 sm:w-24 sm:h-24 md:w-28 md:h-28 rounded-full"
                onClickImage={() => setPreviewOpen(true)}
                showEditIcon
            >
                {/* Image */}
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
                name="name"
                defaultValue={profile?.name}
                placeholder="Enter your full name"
            />
        
            <Input
                className="p-2 sm:p-3 text-sm sm:text-base"
                label="Email"
                type="text"
                name="email"
                disabled
                defaultValue={profile?.email}
            />
        
            <Input
                className="p-2 sm:p-3 text-sm sm:text-base"
                label="Username"
                type="text"
                name="username"
                defaultValue={profile?.username}
                placeholder="Username"
            />
                            
            <div className="grid grid-cols-2 gap-6">
                <Input
                    className="p-2 sm:p-3 text-sm sm:text-base"
                    label="Phone number"
                    type="text"
                    name="phonenumber"
                    defaultValue={profile?.phone_number}
                    placeholder="+1 (555) 000-0000"
                />
        
                <Input
                    className="p-2 sm:p-3 text-sm sm:text-base"
                    label="Country"
                    type="text"
                    name="country"
                    disabled
                    defaultValue={profile?.country}
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
                name="address" 
                defaultValue={profile?.address}
                placeholder="Enter your residential address"
                />
            </div>
        
            <div className="flex flex-col sm:flex-row 
                            gap-6 sm:gap-6 justify-end">
                <Button className="rounded-md px-4 text-purple-500 bg-white hover:bg-[#FAFAFA] ">
                    Back To Home
                </Button>
        
                <Button className="bg-[#312E81] hover:bg-[#432E81] rounded-md px-4">
                    Save Changes
                </Button>
            </div>
        </form>        
    )
}