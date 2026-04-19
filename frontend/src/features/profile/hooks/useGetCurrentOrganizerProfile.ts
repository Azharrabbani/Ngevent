import { useEffect, useState } from "react";
import type { OrganizerResponse } from "../types/profileResponse";
import { GetCurrentOrganizerProfileApi } from "../api/profileApi";

export const useGetCurrentOrganizerProfile = () => {
    const [loading, setLoading] = useState(false)
    const [profile, setProfile] = useState<OrganizerResponse | null>(null)

    const fethProfile = async() => {
        try {
            setLoading(true);        
            const res = await GetCurrentOrganizerProfileApi();
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