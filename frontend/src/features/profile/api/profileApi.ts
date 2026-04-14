import { api } from "../../../lib/api";
import type { successResponse } from "../../../types/apiResponse";
import type { CreateAttendeeProfileReq } from "../types/profileRequest";

export const createAttendeeProfileApi = async (payload: CreateAttendeeProfileReq) => {
    const formData = new FormData();

    formData.append("photo", payload.photo);
    formData.append("name", payload.name);
    formData.append("username", payload.username || "");
    formData.append("phonenumber", payload.phonenumber);
    formData.append("iso", payload.iso);
    formData.append("address", payload.address || "");

    const res = await api.post<successResponse<string>>(
        "/attendee/",
        formData
    );

    return res.data;
};
export const CheckProfile = async () => {
    const res = await api.get<successResponse<boolean>>("/attendee/check-profile");
    return res.data;
}
