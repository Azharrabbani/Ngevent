export interface CreateEventReq {
    id?: string;
    name: string;
    user_id: string;
    description: string;
    categories: number[];
    tickets: ticketsReq[];
    date: number;
    address: eventAddress;
    status: string;
};

interface ticketsReq {
    id?: string;
    name: string;
    price: string;
    quantity: number;
    ticketType: string;
};

interface eventAddress {
    detail_address: string;
    lat: string;
    long: string;
};

