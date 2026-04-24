export const attendeeKeys = {
    all: ["attendees"] as const,
    lists: () => [...attendeeKeys.all, "list"] as const,
    list: (params: any) => [...attendeeKeys.lists(), params] as const,
    details: () => [...attendeeKeys.all, "detail"] as const,
    detail: (id: string) => [...attendeeKeys.details(), id] as const,
    me: () => [...attendeeKeys.all, "me"] as const,
};