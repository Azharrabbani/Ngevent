export const attendeeKeys = {
    all: ["attendees"] as const,
    lists: () => [...attendeeKeys.all, "list"] as const,
    list: (params: any) => [...attendeeKeys.lists(), params] as const,
    details: () => [...attendeeKeys.all, "detail"] as const,
    detail: (id: string) => [...attendeeKeys.details(), id] as const,
    me: () => [...attendeeKeys.all, "me"] as const,
};

export const organizerKeys = {
    all: ["organizers"] as const,
    lists: () => [...organizerKeys.all, "list"] as const,
    list: (params: any) => [...organizerKeys.lists(), params] as const,
    details: () => [...organizerKeys.all, "detail"] as const,
    detail: (id: string) => [...organizerKeys.details(), id] as const,
    me: () => [...organizerKeys.all, "me"] as const,
};

export const organizerUpdateKeys = {
    all: ["organizer-update"] as const,
    detail: (id: string) => [...organizerUpdateKeys.all, id] as const,
};