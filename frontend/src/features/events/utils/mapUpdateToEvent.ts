// src/features/events/utils/mapUpdateToEvent.ts
import type { EventsResponse, UpdateEventResponse } from "../types/eventResponse";

export function mapUpdateToEventResponse(
    update: UpdateEventResponse
): EventsResponse {
    return {
        id: update.event_id,
        update_request_id: update.id,

        eo_profile: update.eo_profile,

        event: {
            banner: update.updated_details.banner,
            name: update.event_title,
            categories: update.updated_categories,
            tickets: [],
            slug: "",
            status: update.updated_details.status,
            description: update.updated_details.description,
        },

        event_address: update.updated_address,

        start_time: update.updated_details.start_time,
        end_time: update.updated_details.end_time,

        created_at: update.created_at,
        updated_at: update.updated_at,
        deleted_at: update.deleted_at,
    };
}