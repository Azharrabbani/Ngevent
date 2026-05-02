import { FaBullhorn } from "react-icons/fa";
import { GoArrowLeft } from "react-icons/go";
import Input from "../../../../../components/input";
import { LuSquarePlus, LuUpload } from "react-icons/lu";
import { FiCalendar } from "react-icons/fi";
import Button from "../../../../../components/Button";
import { useNavigate } from "react-router-dom";

import { useEffect, useState, type Dispatch, type SetStateAction } from "react";
import { toggleItem } from "../../../../../utils/toggleItem";
import SelectItems from "../../../../../components/selectItems";
import toast from "react-hot-toast";
import DatePicker from "react-datepicker";
import type { categoriesResp } from "../../../../categories/types/categoryResponse";
import type { locationResp } from "../../../types/locationResponse";
import MapPicker from "../../../../../components/map";
import TicketForm from "./ticketForm";
import { useForm } from "react-hook-form";
import { useCreateEvent } from "../../hooks/useCreateEvent";
import { converDate } from "../../../../../utils/dateConverter";

type TicketType = "regular" | "premium" | "vip";

type Ticket = {
  tempID: string;
  name: string;
  quantity: string;
  price: string;
  type: TicketType
};

interface Props {
    categories: categoriesResp[] | undefined;
    categoriesLoading: boolean;
    selectedCategories: number[];
    setSelectedCategories: Dispatch<SetStateAction<number[]>>;
    searchQuery: string;
    setSearchQuery: Dispatch<SetStateAction<string>>;
    locations: locationResp[] | undefined;
    locationLoading: boolean;
    selectedLocation: locationResp | undefined;
    setSelectedLocation: Dispatch<SetStateAction<locationResp | undefined>>;
}

export default function EventForm({
    categories,
    categoriesLoading,
    selectedCategories,
    setSelectedCategories,
    searchQuery,
    setSearchQuery,
    locations,
    locationLoading,
    selectedLocation,
    setSelectedLocation,
}: Props) {
    type formValues = {
        name: string;
        detail_address: string;
        description: string;
    }

    const {
        register,
        handleSubmit,
        formState: {errors}
    } = useForm<formValues>();

    const navigate = useNavigate();

    const [isCategoryOpen, setIsCategoryOpen] = useState(false);
    const [tickets, setTickets] = useState<Ticket[]>([
        {
          tempID: crypto.randomUUID(),
          name: "",
          quantity: "",
          price: "",
          type: "regular" as TicketType,
        },
    ]);

    const [banner, setBanner] = useState<File | null>(null);
    const [bannerPreview, setBannerPreview] = useState<string>("");
    const [selectedDate, setSelectedDate] = useState<Date | null>(null);
    const [position, setPosition] = useState<[number, number]>([0,0]);
    const [showDropdown, setShowDropdown] = useState<boolean>(false);


    const createEventMutation = useCreateEvent();

    const toggleCategory = (id: number) => {
        setSelectedCategories((prev) => toggleItem(prev, id));
    };

    const addTicket = () => {
        setTickets((prev) => [
            ...prev,
            {
                tempID: crypto.randomUUID(),
                name: "",
                quantity: "",
                price: "",
                type: "regular" as TicketType,
            },
        ]);
    };

    const updateTicket = (
        tempID: string,
        field: keyof Ticket,
        value: string
    ) => {
        setTickets((prev) => 
            prev.map((ticket) => 
                ticket.tempID === tempID ? { ...ticket, [field]: value} : ticket
        ));
    };

    const removeTicket = (id: string) => {
        setTickets((prev) => prev.filter((ticket) => ticket.tempID !== id));
    };

    const handleBannerUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];

        if (!file) return;

        const allowedTypes = ["image/jpeg", "image/jpg", "image/png"];

        if (!allowedTypes.includes(file.type)) {
            toast.error("Only JPG, JPEG, PNG files are allowed");
            return;
        }

        setBanner(file);
        setBannerPreview(URL.createObjectURL(file));
    };

    const submitEvent = async (
        eventStatus: "draft" | "pending",
        data: formValues
    ) => {
        if (eventStatus === "pending" && !banner) {
            toast.error("Banner required to publish event");
            return;
        }

        if (selectedCategories.length == 0) {
            toast.error("Please select a category");
            return;
        }

        if (!selectedDate) {
            toast.error("Date is required");
            return;
        }

        if (position[0] === 0 || position[1] === 0) {
            toast.error("Please input and select the correct addres");
            return;
        }

        for (const ticket of tickets) {
            if (!ticket.name || !ticket.quantity || !ticket.type) {
                toast.error("All tickets field required");
                return;
            }
        }

        const formattedTickets = tickets.map((ticket) => ({
            name: ticket.name,
            quantity: Number(ticket.quantity),
            price: ticket.price,
            ticket_type: ticket.type,
        }));

        const payload = {
            ...data,
            categories: selectedCategories,
            date: converDate(selectedDate),
            tickets: formattedTickets,
            status: eventStatus,
            address: {
                detail_address: data.detail_address,
                lat: position[0].toString(),
                long: position[1].toString(),
            }
        };

        try {
            await createEventMutation.mutateAsync({
                payload,
                banner,
            });

            navigate("/organizer/dashboard");
        } catch (err) {
            console.error(err);
            toast.error("Failed to create event");
        }
    };

     useEffect(() => {
        return () => {
            if (bannerPreview) {
                URL.revokeObjectURL(bannerPreview);
            }
        };
    }, [bannerPreview]);
    
    return (
        <div className="space-y-10 py-15">
            <div className="flex flex-col sm:flex-row gap-4 sm:gap-8">
                <FaBullhorn className="hidden sm:block" size={45}/>
                <div>
                    <span className="flex justify-center sm:justify-start md:justify-start items-center gap-2">
                        <GoArrowLeft 
                            className="cursor-pointer hover:-translate-x-1 transition duration-180"
                            onClick={() => navigate("/organizer/dashboard")}
                            size={28}/>
                        <h1 className="text-2xl sm:text-3xl text-[#1E293B] font-bold">Create Event</h1>
                    </span>
                    <h2 className="text-sm sm:text-lg text-gray-500">Complete the information to add a new event</h2>
                </div>
            </div>

            <div className="bg-white p-4 sm:p-6 lg:p-8 w-full rounded-xl flex flex-col gap-5">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4 sm:gap-5 mt-1">
                    <Input
                        className="p-2 sm:p-3 text-sm sm:text-base w-full"
                        label="Name"
                        labelStyle="text-[#003B95]"
                        type="text"
                        {...register(
                            "name",
                            {required: "Name is required"}
                        )}
                        error={errors.name?.message}
                        placeholder="e.g. Global Security Summit"
                    />

                    <SelectItems
                        selectedItems={selectedCategories}
                        toogleItem={toggleCategory}
                        isItemOpen={isCategoryOpen}
                        itemsLoading={categoriesLoading}
                        setIsItemOpen={setIsCategoryOpen}
                        items={categories}
                        placeholder="Select categories"
                    />

                    <div className="w-full">
                        <label className="text-[#003B95] text-xs sm:text-sm font-medium mb-1 block">
                            Date
                        </label>

                        <div className="relative z-20">
                            <DatePicker
                                selected={selectedDate}
                                onChange={(date: SetStateAction<Date | null>) => setSelectedDate(date)}
                                dateFormat="dd/MM/yyyy"
                                placeholderText="dd/mm/yyyy"
                                className="
                                    w-full
                                    p-2 sm:p-3
                                    pr-10
                                    text-sm sm:text-base
                                    bg-gray-200
                                    rounded-xl
                                    outline-none
                                "
                                wrapperClassName="w-full"
                                minDate={new Date()}
                            />

                            <FiCalendar
                                className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-500 pointer-events-none"
                                size={18}
                            />
                        </div>
                    </div>
                    
                        
                    <div className="relative w-full">
                        <Input
                            className="p-2 sm:p-3 text-sm sm:text-base"
                            label="Address"
                            labelStyle="text-[#003B95]"
                            type="text"
                            placeholder="Search location..."
                            value={searchQuery}
                            onChange={(e) => {
                                setSearchQuery(e.target.value)
                                setShowDropdown(!showDropdown);
                            }}
                        />

                        {searchQuery && showDropdown && (
                            <div className="absolute z-50 w-full bg-white border mt-1 rounded-xl shadow-lg max-h-60 overflow-y-auto">
                                {locationLoading ? (
                                    <p className="p-3 text-sm text-gray-500">Searching...</p>
                                ) : !locations?.length ? (
                                    <p className="p-3 text-sm text-gray-500">No results</p>
                                ) : (
                                    locations.map((loc) => (
                                        <button
                                            key={loc.display_name}
                                            type="button"
                                            onClick={() => {
                                                setSelectedLocation(loc);
                                                setSearchQuery(loc.display_name);
                                                setPosition([
                                                    parseFloat(loc.lat),
                                                    parseFloat(loc.lon)
                                                ]);
                                                setShowDropdown(false);
                                                
                                            }}
                                            className="w-full text-left px-4 py-2 hover:bg-gray-100 text-sm"
                                        >
                                            {loc.display_name}
                                        </button>
                                    ))
                                )}
                            </div>
                        )}
                    </div>
                </div>

                <MapPicker
                    position={position}
                    selectedLocation={selectedLocation}
                />
                
                <div>
                    <label 
                    className="text-xs sm:text-sm text-[#003B95] font-medium mb-1 block"
                    htmlFor="">
                        Address detail
                    </label>
                    <textarea
                    rows={3}
                    className={`w-full p-2 rounded-xl bg-gray-200 outline-none resize-none`}                         
                    placeholder="Building name, floor, room number, etc."
                    {...register(
                        "detail_address",
                        {required: "Detail of the address required"}
                    )}
                    />
                    {errors.detail_address && (
                        <p className="text-red-500 text-sm">{errors.detail_address.message}</p>
                    )}
                </div>

                <div>
                    <label 
                    className="text-xs text-[#003B95] sm:text-sm font-medium mb-1 block"
                    htmlFor="">
                        Description
                    </label>
                    <textarea
                    rows={3}
                    className={`w-full p-2 rounded-xl bg-gray-200 outline-none resize-none`}                         
                    placeholder="Event description"
                    {...register(
                        "description",
                        {required: "Description required"}
                    )}
                    />
                    {errors.description && (
                        <p className="text-red-500 text-sm">{errors.description.message}</p>
                    )}
                    
                </div>

                {/* Tickets */}
                <div className="bg-[#E2E6EC]/30 p-6 space-y-4 rounded-xl">
                    <div className="flex justify-between items-center">
                        <h2 className="text-[#003B95] font-medium">Tickets</h2>
                        <LuSquarePlus
                            onClick={addTicket} 
                            className="cursor-pointer hover:-translate-y-1 transition duration-150
                                        text-gray-700" 
                            size={20}/>
                    </div>

                    {tickets.map((ticket) => (
                        <TicketForm
                            ticket={ticket}
                            tickets={tickets}
                            updateTicketName={(e) => updateTicket(ticket.tempID, "name", e.target.value)}
                            updateTicketQuanty={(e) => updateTicket(ticket.tempID, "quantity", e.target.value)}
                            updateTicketPrice={(e) => updateTicket(ticket.tempID, "price", e.target.value)}
                            updateTicketType={(value) => updateTicket(ticket.tempID, "type", value)}
                            removeTicket={() => removeTicket(ticket.tempID)}
                        />
                    ))}

                </div>

                {/* Banner */}
                <div className="flex flex-col gap-2 w-full">
                    <p className="text-[#003B95] font-medium">
                        Banner
                    </p>

                    <input
                        id="banner-upload"
                        type="file"
                        accept=".jpg,.jpeg,.png"
                        className="hidden"
                        onChange={handleBannerUpload}
                    />

                    <label
                        htmlFor="banner-upload"
                        className="w-full h-40 sm:h-48 lg:h-56
                                rounded-xl border-2 border-dashed border-gray-300
                                bg-[#E2E6EC]/30
                                flex flex-col gap-2 items-center justify-center
                                cursor-pointer hover:bg-gray-300 transition
                                overflow-hidden"
                    >
                        {bannerPreview ? (
                            <img
                                src={bannerPreview}
                                alt="Banner preview"
                                className="w-full h-full object-cover"
                            />
                        ) : (
                            <>
                                <LuUpload
                                    className="text-gray-500"
                                    size={30}
                                />

                                <p className="text-sm text-gray-600">
                                    Choose a file jpg, jpeg, png
                                </p>
                            </>
                        )}
                    </label>

                    {banner && (
                        <p className="text-sm text-gray-600">
                            {banner.name}
                        </p>
                    )}
                </div>

                <div className="flex flex-col sm:flex-row gap-4 justify-center mt-8 sm:mt-10">
                    <Button
                    type="button" 
                    onClick={handleSubmit((data) => submitEvent("draft", data))}
                    disabled={createEventMutation.isPending}
                    className="w-full sm:w-auto rounded-md px-10 py-3 text-gray-800 font-semibold bg-gray-300 hover:bg-[#c7c7c7]"
                    >
                        {createEventMutation.isPending ? "Loading..." : "Draft"}
                    </Button>
            
                    <Button 
                    type="button"
                    onClick={handleSubmit((data) => submitEvent("pending", data))}
                    disabled={createEventMutation.isPending}
                    className="w-full sm:w-auto rounded-md px-10 py-3 text-white font-semibold bg-[#003B95] hover:bg-[#004ec2]">
                        {createEventMutation.isPending ? "Loading..." : "Create"}
                    </Button>
                </div>
            </div>
        </div>
    )
}