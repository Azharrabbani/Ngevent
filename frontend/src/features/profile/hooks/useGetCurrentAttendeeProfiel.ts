import { useEffect, useState } from "react"
import { GetCurrentAttendeeProfileApi } from "../api/profileApi"
import type { AttendeeResponse } from "../types/profileResponse"

export const useGetCurrentAttendeeProfile = () => {
    const [loading, setLoading] = useState(false)
    const [profile, setProfile] = useState<AttendeeResponse | null>(null)

    const fethProfile = async() => {
        try {
            setLoading(true);        
            const res = await GetCurrentAttendeeProfileApi();
            setProfile(res.data)
        } catch(err: any) {
            console.log(err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fethProfile()
    }, []);
    
    return {profile, loading};
}