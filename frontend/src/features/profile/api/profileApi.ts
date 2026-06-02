import { api } from "../../../lib/api";
import type { PaginatedResponse, successResponse } from "../../../types/apiResponse";
import type { UserFilterRequest } from "../types/userRequest";
import type { UserResponse } from "../types/userResponse";
import type { CreateAttendeeProfileReq, CreateOrganizerProfileReq, UpdatePhotoReq, UpdateAttendeeProfileReq, UpdateOrganizerProfileReq, FilterAttendeeReq, FilterOrganizerReq, validateOrganizerReq, rejectOrganizerReq } from "../types/profileRequest";
import type { AttendeeResponse, OrganizerResponse, OrganizerUpdateResponse } from "../types/profileResponse";

export const CheckProfile = async () => {
    const res = await api.get<successResponse<boolean>>("/attendee/check-profile");
    return res.data;
};


// Attendee profile api
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

export const GetAttendeesProfileApi = async (params: FilterAttendeeReq) => {
    const res = await api.get<PaginatedResponse<AttendeeResponse>>("/attendee", {
        params: {
            filter: params.filter,
            page: params.pagination?.page,
            limit: params.pagination?.limit,
            sort: params.pagination?.sort,
        },
    });

    return res.data;
}

export const GetAttendeeDetailProfileApi = async (id: string) => {
    const res = await api.get<successResponse<AttendeeResponse>>(`/attendee/${id}`);
    return res.data;
}

export const GetAttendeeProfilePhotoApi = async (payload: string) => {
    const res = await api.get(`/attendee/photo/${payload}`);
    return res.data;
}

export const GetCurrentAttendeeProfileApi = async () => {
    const res = await api.get<successResponse<AttendeeResponse>>("/attendee");
    return res.data;
}

export const updateAttendeePhotoApi = async (payload: UpdatePhotoReq) => {
    const formData = new FormData();

    formData.append("photo", payload.photo);

    const res = await api.put<successResponse<string>>(
        "/attendee/photo",
        formData
    );

    return res.data;
}

export const updateAttendeeProfile = async (payload: UpdateAttendeeProfileReq) => {
    const res = await api.put<successResponse<string>>("/attendee", payload);
    return res.data;
}

// Organizer profile api
export const CreateOrganizerProfileApi = async (payload: CreateOrganizerProfileReq) => {
    const formData = new FormData();

    formData.append("photo", payload.photo)
    formData.append("name", payload.name);
    formData.append("phonenumber", payload.phonenumber);
    formData.append("iso", payload.iso);
    formData.append("address", payload.address || "");
    formData.append("nib_number", payload.nib);
    formData.append("npwp_number", payload.npwp);
    formData.append("npwp_file", payload.npwpFile);
    formData.append("nib_file", payload.nibFile);
    formData.append("description", payload.description);

    const res = await api.post<successResponse<string>>(
        "/organizer",
        formData
    );

    return res.data;
};

export const GetOrganizersProfileApi = async (params: FilterOrganizerReq) => {
    const res = await api.get<PaginatedResponse<OrganizerResponse>>("/organizer/profiles", {
        params: {
            filter: params.filter,
            status: params.status,
            page: params.pagination?.page,
            limit: params.pagination?.limit,
            sort: params.pagination?.sort,
        },
    });

    return res.data;
}

export const GetOrganizerDetailProfileApi = async (id: string) => {
    const res = await api.get<successResponse<OrganizerResponse>>(`/organizer/public/${id}`);
    return res.data;
}

export const GetCurrentOrganizerProfileApi = async () => {
    const res = await api.get<successResponse<OrganizerResponse>>("/organizer/me");
    return res.data;
}

export const UpdateOrganizerPhotoApi = async (payload: UpdatePhotoReq) => {
    const formData = new FormData()

    formData.append("photo", payload.photo);

    const res = await api.put<successResponse<string>>(
        "/organizer/photo",
        formData,
    );

    return res.data;
}

export const UpdateOrganizerProfileApi = async (payload: UpdateOrganizerProfileReq) => {
    const formData = new FormData();

    formData.append("name", payload.name);
    formData.append("phonenumber", payload.phonenumber);
    formData.append("iso", payload.iso);
    formData.append("address", payload.address || "");
    formData.append("nib_number", payload.nib);
    formData.append("npwp_number", payload.npwp);

    if (payload.npwpFile) {
        formData.append("npwp_file", payload.npwpFile);
    }

    if (payload.nibFile) {
        formData.append("nib_file", payload.nibFile);
    }

    formData.append("description", payload.description);
    formData.append("email", payload.email);
    formData.append("instagram", payload.instagram);

    const res = await api.put<successResponse<string>>(
        "/organizer",
        formData,
    );

    return res.data;
}

export const GetUpdateOrganizerReqApi = async (id: string) => {
    const res = await api.get<successResponse<OrganizerUpdateResponse>>(`/staging-organizer/update/${id}`);
    return res.data;
}

export const ApproveOrganizerApi = async (id: string) => {
    const res = await api.put<successResponse<string>>(`/organizer/approve/${id}`);
    return res.data;
}

export const RejectOrganizerApi = async (id: string, payload: rejectOrganizerReq) => {
    const res = await api.put<successResponse<string>>(`/organizer/reject/${id}`, payload);
    return res.data;
}

export const ValidateOrganizerUpdateApi = async (id: string | null, payload: validateOrganizerReq) => {
    const res = await api.put<successResponse<string>>(`/staging-organizer/${id}`, payload);
    return res.data;
}

export const CloseOrganizerAccountApi = async () => {
    const res = await api.delete<successResponse<string>>("/organizer/close-account");
    return res.data;
};;

export const listUsersApi = async (params: UserFilterRequest) => {
    const res = await api.get<PaginatedResponse<UserResponse>>(
        "/user",
        {
            params: {
                role: params.role,
                email: params.email,
                page: params.pagination?.page,
                limit: params.pagination?.limit,
                sort: params.pagination?.sort,
            }
        }

    )

    return res.data;
}
