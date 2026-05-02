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
            date: params.date,
            location: params.location,
            page: params.pagination?.page,
            limit: params.pagination?.limit,
            sort: params.pagination?.sort,
        },
    });

    return res.data;
};

export const CreateEventApi = async (payload: CreateEventReq, banner: File | null) => {
    const formData = new FormData();

    formData.append("data", JSON.stringify(payload));

    if (banner) {
        formData.append("banner", banner);
    }

    const res = await api.post<successResponse<string>>("/event", formData);
    return res.data;
}