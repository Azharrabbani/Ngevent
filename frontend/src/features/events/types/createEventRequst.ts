export interface CreateEventReq {
    id?: string;
    name: string;
    user_id: string;
    description: string;
    categories: number[];
    start_time: number;
    end_time: number;
    address: eventAddress;
    status: string;
};

interface eventAddress {
    detail_address: string;
    lat: string;
    long: string;
    display_name?: string;
    city?: string;
    country?: string;
};

