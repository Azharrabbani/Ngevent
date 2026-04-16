import React, { useEffect, useState } from "react";
import Input from "../../../components/input";
import Button from "../../../components/Button";
import type { CreateAttendeeProfileReq } from "../types/profileRequest";
import { GetIsoFromPhoneNumber } from "../../../utils/phoneNumber";

interface Props {
    onSubmit: (payload: CreateAttendeeProfileReq) => void
    loading: boolean
    errors: Record<string, string>
};

export default function AttendeeProfileForm({onSubmit, loading, errors}: Props) {
    const [preview, setPreview] = useState<string | null>(null);
    const [photo, setPhoto] = useState<File | null>(null);
    const [name, setName] = useState<string>("");
    const [username, setUsername] = useState<string>("");
    const [phonenumber, setPhonenumber] = useState<string>("");
    const [address, setAddress] = useState<string>("");

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault()

        if (!photo) {
            alert("Photo is required");
            return;
        }

        const iso = GetIsoFromPhoneNumber(phonenumber)

        if (!iso) {
            alert("Invalid phone number");
            return;
        }

        const payload: CreateAttendeeProfileReq = {
            photo: photo,
            name: name,
            username: username,
            phonenumber: phonenumber,
            iso: iso,
            address: address
        };
        onSubmit(payload)
    }    

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (file) {
            setPreview(URL.createObjectURL(file));
            setPhoto(file)
        }
        
    };

    useEffect(() => {
        return () => {
            if (preview) {
                URL.revokeObjectURL(preview);
            }
        };
    }, [preview]);

    return(
        <form 
            className="flex flex-col gap-4 sm:gap-5 md:gap-6"
            onSubmit={handleSubmit}>
                <div className="flex flex-col justify-center items-center gap-3">
                    <input
                        id="photo-upload"
                        type="file"
                        accept=".jpg, .jpeg, .png"
                        name="photo"
                        className="hidden"
                        onChange={handleFileChange}
                        
                    />
        
                    <label
                        htmlFor="photo-upload"
                        className="
                            w-20 h-20 sm:w-24 sm:h-24 md:w-28 md:h-28
                            rounded-full bg-gray-200
                            flex items-center justify-center
                            cursor-pointer hover:bg-gray-300
                            overflow-hidden transition
                        "
                    >
                        {preview ? (
                             <img src={preview} className="w-full h-full object-cover cursor-pointer" />
                        ): (
                            <svg xmlns="http://www.w3.org/2000/svg" width="25" height="25" fill="currentColor" className="bi bi-plus-square" viewBox="0 0 16 16">
                            <path d="M14 1a1 1 0 0 1 1 1v12a1 1 0 0 1-1 1H2a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1zM2 0a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V2a2 2 0 0 0-2-2z"/>
                            <path d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4"/>
                            </svg>
                        )}
                    </label>
        
                        <label className="text-xs text-center tracking-widest">
                            UPLOAD PHOTO
                        </label>
                </div>
        
                <Input
                    className="p-2 sm:p-3 text-sm sm:text-base"
                    label="Name"
                    type="text"
                    name="name"
                    placeholder="Enter your full name"
                    onChange={(e) => setName(e.target.value)}
                    error={errors.name}
                />
                <Input
                    className="p-2 sm:p-3 text-sm sm:text-base"
                    label="Username"
                    type="text"
                    name="username"
                    placeholder="Username"
                    onChange={(e) => setUsername(e.target.value)}
                />
                <Input
                    className="p-2 sm:p-3 text-sm sm:text-base"
                    label="Phone number"
                    type="text"
                    name="phonenumber"
                    placeholder="+1 (555) 000-0000"
                    onChange={(e) => setPhonenumber(e.target.value)}
                    error={errors.phonenumber}
                />
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
                    placeholder="Enter your residential address"
                    onChange={(e) => setAddress(e.target.value)}
                    />
                </div>

                <Button
                    disabled={loading}
                    type="submit"
                >
                    {loading ? "Loading..." : "Continue"}
                </Button>
        </form>
    )
}