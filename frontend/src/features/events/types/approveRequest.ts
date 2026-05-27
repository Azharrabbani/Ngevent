export interface ReviewEventPayload {
    id: string;
    status: "active" | "rejected";
    reason?: string;
};

export interface ReviewUpdatedEventPayload {
    id: string;
    status: "approved" | "rejected";
    reason?: string;
}
