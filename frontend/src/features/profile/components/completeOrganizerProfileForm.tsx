import React, { useEffect, useState } from "react";
import Input from "../../../components/input";
import Button from "../../../components/Button";
import type { CreateOrganizerProfileReq } from "../types/profileRequest";
import { GetIsoFromPhoneNumber } from "../../../utils/phoneNumber";
import UploadPhoto from "../../../components/uploadPhoto";

interface Props {
    onSubmit: (payload: CreateOrganizerProfileReq) => void
    loading: boolean
    errors: Record<string, string>
};

export default function CompleteOrganizerProfileForm({onSubmit, loading, errors}: Props) {
    const [preview, setPreview] = useState<string | null>(null);
    const [photo, setPhoto] = useState<File | null>(null);
    const [name, setName] = useState<string>("");
    const [description, setDescription] = useState<string>("");
    const [phonenumber, setPhonenumber] = useState<string>("");
    const [nibNumber, setNibNumber] = useState<string>("");
    const [nib, setNIB] = useState<File | null>(null);
    const [npwpNumber, setNpwpNumber] = useState<string>("");
    const [npwp, setNPWP] = useState<File | null>(null);
    const [address, setAddress] = useState<string>("");

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault()

        if (!photo) {
            alert("Photo is required");
            return;
        }

        if (!nib || !npwp) {
            alert("nib and nwpwp must be upload")
            return;
        }

        const iso = GetIsoFromPhoneNumber(phonenumber)

        if (!iso) {
            alert("Invalid phone number");
            return;
        }

        const payload: CreateOrganizerProfileReq = {
            photo: photo,
            name: name,
            phonenumber: phonenumber,
            iso: iso,
            address: address,
            description: description,
            nib: nibNumber,
            nibFile: nib,
            npwp: npwpNumber,
            npwpFile: npwp,
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

    const handleNpwpChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (!file) return;

        if (file.type !== "application/pdf") {
            alert("NPWP must be PDF");
            return;
        }

        setNPWP(file);   
    }

    const handleNibChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (!file) return;

        if (file.type !== "application/pdf") {
            alert("NIB must be PDF");
            return;
        }

        setNIB(file);   
    }


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
                <UploadPhoto
                    className="mt-5 relative z-10 w-20 h-20 sm:w-24 sm:h-24 md:w-28 md:h-28 rounded-full"
                    onChange={handleFileChange}
                >
                    {preview ? (
                             <img src={preview} className="w-full h-full object-cover cursor-pointer rounded-full" />
                        ): (
                            <svg xmlns="http://www.w3.org/2000/svg" width="21" height="21" fill="currentColor" className="bi bi-plus-square" viewBox="0 0 16 16">
                                <path d="M14 1a1 1 0 0 1 1 1v12a1 1 0 0 1-1 1H2a1 1 0 0 1-1-1V2a1 1 0 0 1 1-1zM2 0a2 2 0 0 0-2 2v12a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V2a2 2 0 0 0-2-2z"/>
                                <path d="M8 4a.5.5 0 0 1 .5.5v3h3a.5.5 0 0 1 0 1h-3v3a.5.5 0 0 1-1 0v-3h-3a.5.5 0 0 1 0-1h3v-3A.5.5 0 0 1 8 4"/>
                            </svg>
                    )}
                    
                </UploadPhoto>
                <label className="text-xs text-center tracking-widest">
                    ORGANIZER AVATAR
                </label>
        
                <Input
                    className="p-2 sm:p-3 text-sm sm:text-base"
                    label="Name"
                    type="text"
                    name="name"
                    placeholder="Enter your full name"
                    onChange={(e) => setName(e.target.value)}
                    error={errors.name}
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
                    name="address" 
                    placeholder="Enter your residential address"
                    onChange={(e) => setDescription(e.target.value)}
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
                    className={`w-full p-2 rounded-xl bg-gray-200 outline-none resize-none ${
                        errors.address ? "border border-red-500" : ""
                    }`}                         
                    name="address" 
                    placeholder="Enter your residential address"
                    onChange={(e) => setAddress(e.target.value)}
                    />
                    {errors.address && <p className="text-red-500">{errors.address}</p>}
                </div>

                <Input
                    className="p-2 sm:p-3 text-sm sm:text-base"
                    label="Phone number"
                    type="text"
                    name="phonenumber"
                    placeholder="+1 (555) 000-0000"
                    onChange={(e) => setPhonenumber(e.target.value)}
                    error={errors.phonenumber}
                />

                <Input
                    className="p-2 sm:p-3 text-sm sm:text-base"
                    label="NPWP number"
                    type="text"
                    name="npwp"
                    placeholder="00-000-000-00"
                    onChange={(e) => setNpwpNumber(e.target.value)}
                    error={errors.npwp}
                />

                <Input
                    className="p-2 sm:p-3 text-sm sm:text-base"
                    label="NIB number"
                    type="text"
                    name="nib"
                    placeholder="00-000-000-00"
                    onChange={(e) => setNibNumber(e.target.value)}
                    error={errors.nib}
                />

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    {/* NPWP */}
                    <div className="flex flex-col items-center gap-2 w-full">
                        <p className="text-xs tracking-widest text-center">
                            NPWP FILE
                        </p>

                        <input
                            id="npwp-upload"
                            type="file"
                            accept="application/pdf"
                            className="hidden"
                            onChange={handleNpwpChange}
                        />

                        <label
                            htmlFor="npwp-upload"
                            className={`
                                    w-full h-28 sm:h-32
                                    rounded-xl border-2 border-dashed border-gray-400
                                    bg-gray-100
                                    flex flex-col gap-2 items-center justify-center
                                    cursor-pointer hover:bg-gray-200
                                    transition
                                    ${errors.npwp && "border border-red-500"}
                                `}
                        >
                            {npwp ? (
                                <>
                                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="red" className="bi bi-filetype-pdf" viewBox="0 0 16 16">
                                    <path fill-rule="evenodd" d="M14 4.5V14a2 2 0 0 1-2 2h-1v-1h1a1 1 0 0 0 1-1V4.5h-2A1.5 1.5 0 0 1 9.5 3V1H4a1 1 0 0 0-1 1v9H2V2a2 2 0 0 1 2-2h5.5zM1.6 11.85H0v3.999h.791v-1.342h.803q.43 0 .732-.173.305-.175.463-.474a1.4 1.4 0 0 0 .161-.677q0-.375-.158-.677a1.2 1.2 0 0 0-.46-.477q-.3-.18-.732-.179m.545 1.333a.8.8 0 0 1-.085.38.57.57 0 0 1-.238.241.8.8 0 0 1-.375.082H.788V12.48h.66q.327 0 .512.181.185.183.185.522m1.217-1.333v3.999h1.46q.602 0 .998-.237a1.45 1.45 0 0 0 .595-.689q.196-.45.196-1.084 0-.63-.196-1.075a1.43 1.43 0 0 0-.589-.68q-.396-.234-1.005-.234zm.791.645h.563q.371 0 .609.152a.9.9 0 0 1 .354.454q.118.302.118.753a2.3 2.3 0 0 1-.068.592 1.1 1.1 0 0 1-.196.422.8.8 0 0 1-.334.252 1.3 1.3 0 0 1-.483.082h-.563zm3.743 1.763v1.591h-.79V11.85h2.548v.653H7.896v1.117h1.606v.638z"/>
                                    </svg>
                                    <p className="text-xs text-gray-600 truncate max-w-full text-center px-2">
                                        {npwp.name}
                                    </p>
                                </>
                            ): (
                                <>
                                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="gray" viewBox="0 0 16 16">
                                        <path d="M.5 9.9a.5.5 0 0 1 .5.5v2.5a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1v-2.5a.5.5 0 0 1 1 0v2.5a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2v-2.5a.5.5 0 0 1 .5-.5"/>
                                        <path d="M7.646 1.146a.5.5 0 0 1 .708 0l3 3a.5.5 0 0 1-.708.708L8.5 2.707V11.5a.5.5 0 0 1-1 0V2.707L5.354 4.854a.5.5 0 1 1-.708-.708z"/>
                                    </svg>   
                                    <p className="text-xs text-gray-600">Click to upload NPWP</p>
                                </>
                            )}
                        </label>
                        {errors.npwp && <p className="text-red-500">{errors.npwp}</p>}
                    </div>

                    {/* NIB */}
                    <div className="flex flex-col items-center gap-2 w-full">
                        <p className="text-xs tracking-widest text-center">
                            NIB FILE
                        </p>

                        <input
                            id="nib-upload"
                            type="file"
                            accept="application/pdf"
                            className="hidden"
                            onChange={handleNibChange}
                        />

                        <label
                            htmlFor="nib-upload"
                            className={`
                                    w-full h-28 sm:h-32
                                    rounded-xl border-2 border-dashed border-gray-400
                                    bg-gray-100
                                    flex flex-col gap-2 items-center justify-center
                                    cursor-pointer hover:bg-gray-200
                                    transition
                                    ${errors.nib && "border border-red-500"}
                                `}
                        >
                            {nib ? (
                                <>
                                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="red" className="bi bi-filetype-pdf" viewBox="0 0 16 16">
                                        <path fill-rule="evenodd" d="M14 4.5V14a2 2 0 0 1-2 2h-1v-1h1a1 1 0 0 0 1-1V4.5h-2A1.5 1.5 0 0 1 9.5 3V1H4a1 1 0 0 0-1 1v9H2V2a2 2 0 0 1 2-2h5.5zM1.6 11.85H0v3.999h.791v-1.342h.803q.43 0 .732-.173.305-.175.463-.474a1.4 1.4 0 0 0 .161-.677q0-.375-.158-.677a1.2 1.2 0 0 0-.46-.477q-.3-.18-.732-.179m.545 1.333a.8.8 0 0 1-.085.38.57.57 0 0 1-.238.241.8.8 0 0 1-.375.082H.788V12.48h.66q.327 0 .512.181.185.183.185.522m1.217-1.333v3.999h1.46q.602 0 .998-.237a1.45 1.45 0 0 0 .595-.689q.196-.45.196-1.084 0-.63-.196-1.075a1.43 1.43 0 0 0-.589-.68q-.396-.234-1.005-.234zm.791.645h.563q.371 0 .609.152a.9.9 0 0 1 .354.454q.118.302.118.753a2.3 2.3 0 0 1-.068.592 1.1 1.1 0 0 1-.196.422.8.8 0 0 1-.334.252 1.3 1.3 0 0 1-.483.082h-.563zm3.743 1.763v1.591h-.79V11.85h2.548v.653H7.896v1.117h1.606v.638z"/>
                                    </svg>
                                    <p className="text-xs text-gray-600 truncate max-w-full text-center px-2">
                                        {nib.name}
                                    </p>
                                </>
                            ) : (
                                <>
                                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="gray" viewBox="0 0 16 16">
                                        <path d="M.5 9.9a.5.5 0 0 1 .5.5v2.5a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1v-2.5a.5.5 0 0 1 1 0v2.5a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2v-2.5a.5.5 0 0 1 .5-.5"/>
                                        <path d="M7.646 1.146a.5.5 0 0 1 .708 0l3 3a.5.5 0 0 1-.708.708L8.5 2.707V11.5a.5.5 0 0 1-1 0V2.707L5.354 4.854a.5.5 0 1 1-.708-.708z"/>
                                    </svg>   
                                    <p className="text-xs text-gray-600">Click to upload NIB</p>
                                </>
                            )}
                        </label>
                        {errors.nib && <p className="text-red-500">{errors.nib}</p>}
                    </div>

                </div>

                <Button
                    className="bg-[#312E81] hover:bg-purple-900"
                    disabled={loading}
                    type="submit"
                >
                    {loading ? "Loading..." : "Continue"}
                </Button>
        </form>
    )
}