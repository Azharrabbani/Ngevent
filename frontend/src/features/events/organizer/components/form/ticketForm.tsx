import { FiX } from "react-icons/fi";
import { IoIosArrowDown } from "react-icons/io";
import Input from "../../../../../components/input";

interface Props {
    ticket: any;
    tickets: any[];
    updateTicketName: (e: any) => void;
    updateTicketQuanty: (e: any) => void;
    updateTicketPrice: (e: any) => void;
    updateTicketType: (value: string) => void;
    removeTicket: () => void;
};


export default function TicketForm({
    ticket, 
    tickets,
    updateTicketName, 
    updateTicketQuanty,
    updateTicketPrice,
    updateTicketType,
    removeTicket,
}: Props) {
    const TICKET_TYPES = [
      { value: "regular", label: "Regular" },
      { value: "premium", label: "Premium" },
      { value: "vip", label: "VIP" },
    ];
    return (
        <div className="relative grid grid-cols-1 md:grid-cols-[1fr_150px_200px_1fr] gap-4 sm:gap-5
                        bg-[#c5c5c5] p-5 rounded-xl"
                        
        >
            <Input
                className="p-2 sm:p-3 text-sm sm:text-base w-full"
                label="Name"
                labelStyle="text-gray-700"
                type="text"
                placeholder="e.g. VIP Pass"
                value={ticket.name}
                onChange={updateTicketName}
            />
            <Input
                className="p-2 sm:p-3 text-sm sm:text-base w-full"
                label="Quantity"
                labelStyle="text-gray-700"
                type="text"
                inputMode="numeric"
                placeholder="0"
                value={ticket.quantity}
                onChange={updateTicketQuanty}
            />
            <Input
                className="p-2 sm:p-3 text-sm sm:text-base w-full"
                label="Price"
                labelStyle="text-gray-700"
                type="text"
                inputMode="numeric"
                placeholder="IDR 0.00"
                value={ticket.price}
                onChange={updateTicketPrice}
            />
            <div className="w-full">
                <label className="text-gray-700 text-sm font-medium mb-1 block">
                    Type
                </label>

                <div className="relative">
                    <select
                    value={ticket.type}
                    onChange={(e) => updateTicketType(e.target.value)}
                    className="
                        w-full
                        p-2 sm:p-3
                        text-sm sm:text-base
                        bg-gray-200
                        rounded-xl
                        outline-none
                        appearance-none
                        pr-10
                        cursor-pointer
                        focus:ring-2 focus:ring-blue-400
                    "
                    >
                        <option value="">Select type</option>
                        {TICKET_TYPES.map((type) => (
                            <option key={type.value} value={type.value}>
                                {type.label}
                            </option>
                        ))}
                    </select>

                    <div 
                        className="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-gray-500"
                    >
                        <IoIosArrowDown/>
                    </div>
                </div>
            </div>

            {tickets.length > 1 && (
                <button
                    type="button"
                    onClick={removeTicket}
                    className="absolute self-end mb-3 top-2 right-2"
                >
                    <FiX
                    className="text-gray-600 hover:text-gray-700"
                    size={23}
                    />
                </button>
            )}
        </div>
    )
}