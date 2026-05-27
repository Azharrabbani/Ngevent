import api from "../../../lib/api"
import type { PaginatedResponse, successResponse } from "../../../types/apiResponse"
import type { CreateEventReq } from "../types/createEventRequst"
import type { locationReq } from "../types/locationRequest"
import type { locationResp } from "../types/locationResponse"
import type { FilterEventsRequest, FilterUpdatedEventsRequest } from "../types/eventRequest"
import type { EventsResponse, UpdateEventResponse } from "../types/eventResponse"

export const SearchLocationApi = async (payload: locationReq) => {
    const res = await api.get<successResponse<locationResp[]>>("/location/", {
        params: {
            query: payload.query
        },
    });

    return res.data;
}

export const GetEventsApi = async (params: FilterEventsRequest) => {
    const res = await api.get<PaginatedResponse<EventsResponse>>("/event", {
        params: {
            search: params.search,
            status: params.status,
            page: params.pagination?.page,
            limit: params.pagination?.limit,
            sort: params.sort,
            date: params.date,
        }
    });

    return res.data;
}

export const GetAllUpdatedEventsApi = async (params: FilterUpdatedEventsRequest) => {
    const res = await api.get<PaginatedResponse<UpdateEventResponse>>("/updated-event", {
        params: {
            title: params.title,
            search: params.search,
            sort: params.sort,
            date: params.date,
            status: params.status,
            page: params.pagination?.page,
            limit: params.pagination?.limit,
        },
    });
    return res.data;
};

export const GetOrganizerEventsApi = async (params: FilterEventsRequest) => {
    const res = await api.get<PaginatedResponse<EventsResponse>>("event/organizer-events", {
        params: {
            title: params.title,
            category: params.category,
            status: params.status,
            start_time: params.start_time,
            location: params.location,
            page: params.pagination?.page,
            limit: params.pagination?.limit,
            sort: params.pagination?.sort,
        },
    });

    return res.data;
};

export const GetEventByID = async (id: string) => {
    const res = await api.get<successResponse<EventsResponse>>(`event/${id}`);
    return res.data;
}

export const GetUpdateByEventIDApi = async (eventID: string, status: string) => {
    const res = await api.get<successResponse<UpdateEventResponse>>(`updated-event/${eventID}`, {
        params: {
            status: status,
        },
    });
    return res.data;
}


export const CreateEventApi = async (payload: CreateEventReq, banner: File | null) => {
    const formData = new FormData();

    formData.append("data", JSON.stringify(payload));

    if (banner) {
        formData.append("banner", banner);
    }

    const res = await api.post<successResponse<string>>("/event", formData);
    return res.data;
}

export const UpdateEventApi = async (id: string, payload: CreateEventReq, banner: File | null) => {
    const formData = new FormData();

    formData.append("data", JSON.stringify(payload));

    if (banner) {
        formData.append("banner", banner);
    }

    const res = await api.put<successResponse<string>>(`/event/${id}`, formData);
    return res.data;
};
export const CancelEventApi = async (id: string) => {
    const res = await api.put<successResponse<string>>(`/event/cancel/${id}`)
    return res.data
}

export const DeleteEventApi = async (id: string) => {
    const res = await api.delete<successResponse<string>>(`/event/${id}`)
    return res.data
}


export const ReviewEventApi = async (id: string, payload: { status: "active" | "rejected"; reason?: string }) => {
    const res = await api.put<successResponse<string>>(`/event/review/${id}`, payload);
    return res.data;
};

export const ReviewUpdatedEventApi = async (id: string, payload: { status: "approved" | "rejected"; reason?: string }) => {
    const res = await api.put<successResponse<string>>(`/updated-event/review/${id}`, payload);
    return res.data;
};