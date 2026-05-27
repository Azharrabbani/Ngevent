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

export const eventsKeys = {
    all: ["events"] as const,
    lists: () => [...eventsKeys.all, "list"] as const,
    list: (params: any) => [...eventsKeys.lists(), params] as const,
    details: () => [...eventsKeys.all, "detail"] as const,
    detail: (id: string) => [...eventsKeys.details(), id] as const,
    me: () => [...eventsKeys.all, "me"] as const,
};

export const updateEventKeys = {
    all: ["update-event"] as const,
    lists: () => [...updateEventKeys.all, "list"] as const,
    list: (params: any) => [...updateEventKeys.lists(), params] as const,
    details: () => [...updateEventKeys.all, "detail"] as const,
    detail: (id: string, status?: string) => [...updateEventKeys.details(), id, status || ""] as const,
    me: () => [...updateEventKeys.all, "me"] as const,
};

export const categoriesKeys = {
    all: ["categories"] as const,
    lists: () => [...categoriesKeys.all, "list"] as const,
    list: (params: any) => [...categoriesKeys.lists(), params] as const,
    details: () => [...categoriesKeys.all, "detail"] as const,
    detail: (id: string) => [...categoriesKeys.details(), id] as const,
    me: () => [...categoriesKeys.all, "me"] as const,
};

export const locationKeys = {
    all: ["locations"] as const,
    lists: () => [...locationKeys.all, "list"] as const,
    list: (params: any) => [...locationKeys.lists(), params] as const,
};