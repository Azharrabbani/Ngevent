import api from "../../../../lib/api"
import type { PaginatedResponse, successResponse } from "../../../../types/apiResponse";
import type { CreateEventReq } from "../../types/createEventRequst";
import type { FilterEventsRequest } from "../../types/organizerRequest"
import type { EventsResponse } from "../../types/organizerResponse"

export const GetEvents = async (params: FilterEventsRequest) => {
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
}

export const CancelEventApi = async (id: string) => {
    const res = await api.put<successResponse<string>>(`/event/cancel/${id}`)
    return res.data
}

export const DeleteEventApi = async (id: string) => {
    const res = await api.delete<successResponse<string>>(`/event/${id}`)
    return res.data
}