import { FaBullhorn } from "react-icons/fa";
import { GoArrowLeft } from "react-icons/go";
import Input from "../../../../../components/input";
import { LuUpload } from "react-icons/lu";
import { FiCalendar } from "react-icons/fi";
import Button from "../../../../../components/Button";
import { useNavigate, useParams } from "react-router-dom";

import { useEffect, useState, type Dispatch, type SetStateAction } from "react";
import { toggleItem } from "../../../../../utils/toggleItem";
import SelectItems from "../../../../../components/selectItems";
import toast from "react-hot-toast";
import DatePicker from "react-datepicker";
import type { categoriesResp } from "../../../../categories/types/categoryResponse";
import type { locationResp } from "../../../types/locationResponse";
import MapPicker from "../../../../../components/map";
import { useForm } from "react-hook-form";
import type { EventsResponse } from "../../../types/eventResponse";
import RichTextEditor from "../../../../../components/richTextEditior";
import { useCreateEvent } from "../../../hooks/useCreateEvent";
import { useUpdateEvent } from "../../../hooks/useUpdateEvent";
import { useCancelEvent } from "../../../hooks/useCancelEvent";
import { useDeleteEvent } from "../../../hooks/useDeleteEvent";


interface Props {
    mode: "create" | "edit";
    eventData?: EventsResponse;
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
    mode,
    eventData,
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
    const { id } = useParams<{ id: string }>();
    const isEditMode = mode === "edit";
    const eventStatus = eventData?.event.status;

    type formValues = {
        name: string;
        detail_address: string;
        description: string;
    };

    const {
        register,
        handleSubmit,
        reset,
        watch,
        setValue,
        formState: { errors },
    } = useForm<formValues>({
        defaultValues: {
            description: ""
        }
    });

    const navigate = useNavigate();

    const [isCategoryOpen, setIsCategoryOpen] = useState(false);

    const [banner, setBanner] = useState<File | null>(null);
    const [bannerPreview, setBannerPreview] = useState<string>("");
    const [selectedDate, setSelectedDate] = useState<Date | null>(null);
    const [startTime, setStartTime] = useState<Date | null>(null)
    const [endTime, setEndTime] = useState<Date | null>(null)
    const [position, setPosition] = useState<[number, number]>([0, 0]);
    const [showDropdown, setShowDropdown] = useState<boolean>(false);

    const createEventMutation = useCreateEvent();
    const updateEventMutation = useUpdateEvent();
    const cancelEventMutation = useCancelEvent();
    const deleteEventMutation = useDeleteEvent();

    const [confirmAction, setConfirmAction] = useState<"cancel" | "delete" | null>(null)

    const isPending = isEditMode
        ? updateEventMutation.isPending
        : createEventMutation.isPending;

    // Pre-populate form when in edit mode
    useEffect(() => {
        if (!isEditMode || !eventData) return;

        // Reset form fields
        reset({
            name: eventData.event.name,
            detail_address: eventData.event_address.detail_address,
            description: eventData.event.description,
        });

        // Set categories
        const categoryIds = eventData.event.categories.map((c) => Number(c.id));
        setSelectedCategories(categoryIds);

        // Set date from unix timestamp
        if (eventData.start_time) {
            const start = new Date(eventData.start_time * 1000)
            const end = new Date(eventData.end_time * 1000)
            setSelectedDate(start)
            setStartTime(start)
            setEndTime(end)
        }

        // Set position from coordinates
        const { lat, lon } = eventData.event_address.coordinates;
        if (lat && lon) {
            setPosition([lat, lon]);
            setSearchQuery(eventData.event_address.address);
        }

        // Set banner preview from existing banner URL
        if (eventData.event.banner) {
            setBannerPreview(eventData.event.banner);
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [isEditMode, eventData]);

    const toggleCategory = (catId: number) => {
        setSelectedCategories((prev) => toggleItem(prev, catId));
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

    const validateForm = (targetStatus: string): boolean => {
        // Banner required when publishing (not draft)
        if (targetStatus !== "draft" && !isEditMode && !banner) {
            toast.error("Banner required to publish event");
            return false;
        }

        if (selectedCategories.length === 0) {
            toast.error("Please select a category");
            return false;
        }

        if (!selectedDate) {
            toast.error("Date is required");
            return false;
        }

        if (!startTime) {
            toast.error("Start time is required")
            return false
        }

        if (!endTime) {
            toast.error("End time is required")
            return false
        }

        if (endTime.getTime() <= startTime.getTime()) {
            toast.error("End time must be after start time")
            return false
        }

        if (position[0] === 0 || position[1] === 0) {
            toast.error("Please input and select the correct address");
            return false;
        }

        return true;
    };

    const buildPayload = (targetStatus: string, data: formValues) => {
        const mergeDateTime = (date: Date, time: Date) => {
            const merged = new Date(date)
            merged.setHours(time.getHours(), time.getMinutes(), 0, 0)
            return Math.floor(merged.getTime() / 1000)
        }
        return {
            ...data,
            categories: selectedCategories,
            start_time: mergeDateTime(selectedDate!, startTime!),
            end_time: mergeDateTime(selectedDate!, endTime!),
            status: targetStatus,
            address: {
                detail_address: data.detail_address,
                lat: position[0].toString(),
                long: position[1].toString(),
            },
        };
    };

    // Create Mode Handler
    const handleCreate = async (targetStatus: "draft" | "pending", data: formValues) => {
        if (!validateForm(targetStatus)) return;
        const payload = buildPayload(targetStatus, data);

        try {
            await createEventMutation.mutateAsync({ payload, banner });
            navigate(-1);
        } catch (err) {
            console.error(err);
        }
    };

    // Update Mode Handler
    const handleUpdate = async (targetStatus: string, data: formValues) => {
        if (!id) return;
        if (!validateForm(targetStatus)) return;
        const payload = buildPayload(targetStatus, data);

        try {
            await updateEventMutation.mutateAsync({ id, payload, banner });
            navigate(-1);
        } catch (err) {
            console.error(err);
        }
    };

    const handleCancel = async () => {
        if (!id) return
        await cancelEventMutation.mutateAsync(id)
        navigate(-1)
    }

    const handleDelete = async () => {
        if (!id) return
        await deleteEventMutation.mutateAsync(id)
        navigate(-1)
    }

    useEffect(() => {
        return () => {
            if (bannerPreview && bannerPreview.startsWith("blob:")) {
                URL.revokeObjectURL(bannerPreview);
            }
        };
    }, [bannerPreview]);

    // Button type depends on mode (create/update) and event status (draft/pending/active)
    const renderButtons = () => {
        if (!isEditMode) {
            // Create mode: Draft + Create
            return (
                <>
                    <Button
                        type="button"
                        onClick={handleSubmit((data) => handleCreate("draft", data))}
                        disabled={isPending}
                        className="w-full sm:w-auto rounded-md px-10 py-3 text-gray-800 font-semibold bg-gray-300 hover:bg-[#c7c7c7]"
                    >
                        {isPending ? "Loading..." : "Draft"}
                    </Button>

                    <Button
                        type="button"
                        onClick={handleSubmit((data) => handleCreate("pending", data))}
                        disabled={isPending}
                        className="w-full sm:w-auto rounded-md px-10 py-3 text-white font-semibold bg-[#003B95] hover:bg-[#004ec2]"
                    >
                        {isPending ? "Loading..." : "Create"}
                    </Button>
                </>
            );
        }

        if (eventStatus === "pending") {
            return null; // No buttons available for pending events
        }

        if (eventStatus === "cancelled") {
            return (
                <>
                    <Button
                        type="button"
                        onClick={handleSubmit((data) => handleUpdate("draft", data))}
                        disabled={isPending}
                        className="w-full sm:w-auto rounded-md px-10 py-3 text-gray-800 font-semibold bg-gray-300 hover:bg-[#c7c7c7]"
                    >
                        {isPending ? "Loading..." : "Save as Draft"}
                    </Button>

                    <Button
                        type="button"
                        onClick={handleSubmit((data) => handleUpdate("pending", data))}
                        disabled={isPending}
                        className="w-full sm:w-auto rounded-md px-10 py-3 text-white font-semibold bg-[#003B95] hover:bg-[#004ec2]"
                    >
                        {isPending ? "Loading..." : "Republish"}
                    </Button>
                </>
            )
        }

        if (eventStatus === "draft") {
            // Edit draft: Update (keep draft) + Publish (pending)
            return (
                <>
                    <Button
                        type="button"
                        onClick={() => setConfirmAction("delete")}
                        disabled={deleteEventMutation.isPending}
                        className="w-full sm:w-auto rounded-md px-10 py-3 text-white font-semibold bg-red-500 hover:bg-red-600"
                    >
                        {deleteEventMutation.isPending ? "Deleting..." : "Delete"}
                    </Button>
                    <Button
                        type="button"
                        onClick={handleSubmit((data) => handleUpdate("draft", data))}
                        disabled={isPending}
                        className="w-full sm:w-auto rounded-md px-10 py-3 text-gray-800 font-semibold bg-gray-300 hover:bg-[#c7c7c7]"
                    >
                        {isPending ? "Loading..." : "Update"}
                    </Button>

                    <Button
                        type="button"
                        onClick={handleSubmit((data) => handleUpdate("pending", data))}
                        disabled={isPending}
                        className="w-full sm:w-auto rounded-md px-10 py-3 text-white font-semibold bg-[#003B95] hover:bg-[#004ec2]"
                    >
                        {isPending ? "Loading..." : "Publish"}
                    </Button>
                </>
            );
        }

        // Edit active event: Update only (goes to admin review)
        return (
            <>
                <Button
                    type="button"
                    onClick={() => setConfirmAction("cancel")}
                    disabled={cancelEventMutation.isPending}
                    className="w-full sm:w-auto rounded-md px-10 py-3 text-white font-semibold bg-red-500 hover:bg-red-600"
                >
                    {cancelEventMutation.isPending ? "Canceling..." : "Cancel Event"}
                </Button>

                <Button
                    type="button"
                    onClick={handleSubmit((data) => handleUpdate("pending", data))}
                    disabled={isPending}
                    className="w-full sm:w-auto rounded-md px-10 py-3 text-white font-semibold bg-[#003B95] hover:bg-[#004ec2]"
                >
                    {isPending ? "Loading..." : "Update"}
                </Button>
            </>
        );
    };

    const pageTitle = isEditMode ? "Edit Event" : "Create Event";
    const pageSubtitle = isEditMode
        ? "Update the information for this event"
        : "Complete the information to add a new event";

    return (
        <div className="w-full max-w-7xl mx-auto space-y-10 px-4 sm:px-6 lg:px-8 py-10">
            <div className="flex flex-col sm:flex-row gap-4 sm:gap-8">
                <FaBullhorn className="hidden sm:block" size={45} />
                <div>
                    <span className="flex justify-center sm:justify-start md:justify-start items-center gap-2">
                        <GoArrowLeft
                            className="cursor-pointer hover:-translate-x-1 transition duration-180"
                            onClick={() => navigate(-1)}
                            size={28}
                        />
                        <h1 className="text-3xl sm:text-4xl lg:text-5xl text-[#1E293B] font-bold">{pageTitle}</h1>
                    </span>
                    <h2 className="text-center sm:text-base sm:text-lg text-gray-500 mt-1">{pageSubtitle}</h2>
                </div>
            </div>

            <div className="bg-white w-full rounded-2xl shadow-sm border border-gray-100 p-5 sm:p-8 lg:p-10 xl:p-12 flex flex-col gap-8">
                <div className="grid grid-cols-1 xl:grid-cols-2 gap-6 lg:gap-8 mt-2">
                    <Input
                        className="w-full rounded-xl p-3 sm:p-4 text-sm sm:text-base"
                        label="Name"
                        labelStyle="text-[#003B95]"
                        type="text"
                        {...register("name", { required: "Name is required" })}
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

                    {selectedDate && (
                        <div className="w-full grid grid-cols-2 gap-4">
                            <div>
                                <label className="text-[#003B95] text-xs sm:text-sm font-medium mb-1 block">
                                    Start Time (WIB)
                                </label>
                                <div className="relative z-20">
                                    <DatePicker
                                        selected={startTime}
                                        onChange={(time: SetStateAction<Date | null>) => setStartTime(time)}
                                        showTimeSelect
                                        showTimeSelectOnly
                                        timeIntervals={15}
                                        timeCaption="Start"
                                        dateFormat="HH:mm"
                                        placeholderText="09:00 WIB"
                                        className="w-full p-2 sm:p-3 text-sm sm:text-base bg-gray-200 rounded-xl outline-none"
                                        wrapperClassName="w-full"
                                    />
                                </div>
                            </div>

                            <div>
                                <label className="text-[#003B95] text-xs sm:text-sm font-medium mb-1 block">
                                    End Time (WIB)
                                </label>
                                <div className="relative z-20">
                                    <DatePicker
                                        selected={endTime}
                                        onChange={(time: SetStateAction<Date | null>) => setEndTime(time)}
                                        showTimeSelect
                                        showTimeSelectOnly
                                        timeIntervals={15}
                                        timeCaption="End"
                                        dateFormat="HH:mm"
                                        placeholderText="17:00 WIB"
                                        minTime={startTime ? new Date(startTime.getTime() + 15 * 60 * 1000) : undefined}
                                        maxTime={new Date(new Date().setHours(23, 45))}
                                        className="w-full p-2 sm:p-3 text-sm sm:text-base bg-gray-200 rounded-xl outline-none"
                                        wrapperClassName="w-full"
                                    />
                                </div>
                            </div>
                        </div>
                    )}

                    <div className="relative w-full">
                        <Input
                            className="w-full rounded-xl p-3 sm:p-4 text-sm sm:text-base"
                            label="Address"
                            labelStyle="text-[#003B95]"
                            type="text"
                            placeholder="Search location..."
                            value={searchQuery}
                            onChange={(e) => {
                                setSearchQuery(e.target.value);
                                setShowDropdown(true);
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
                                                    parseFloat(loc.lon),
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

                <MapPicker position={position} selectedLocation={selectedLocation} />

                <div>
                    <label
                        className="text-xs sm:text-sm text-[#003B95] font-medium mb-1 block"
                        htmlFor="detail_address"
                    >
                        Address detail
                    </label>
                    <textarea
                        rows={5}
                        className="w-full rounded-xl bg-gray-200 p-4 text-sm sm:text-base outline-none resize-none"
                        placeholder="Building name, floor, room number, etc."
                        {...register("detail_address", {
                            required: "Detail of the address required",
                        })}
                    />
                    {errors.detail_address && (
                        <p className="text-red-500 text-sm">{errors.detail_address.message}</p>
                    )}
                </div>

                <div>
                    <label
                        className="text-xs text-[#003B95] sm:text-sm font-medium mb-1 block"
                        htmlFor="description"
                    >
                        Description
                    </label>
                    <RichTextEditor
                        value={watch("description") ?? ""}
                        onChange={(val) => setValue("description", val, { shouldDirty: true })}
                        error={errors.description?.message}
                    />
                </div>

                {/* Banner */}
                <div className="flex flex-col gap-2 w-full">
                    <p className="text-[#003B95] font-medium">Banner</p>

                    <input
                        id="banner-upload"
                        type="file"
                        accept=".jpg,.jpeg,.png"
                        className="hidden"
                        onChange={handleBannerUpload}
                    />

                    <label
                        htmlFor="banner-upload"
                        className="w-full h-52 sm:h-64 lg:h-80 xl:h-96
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
                                <LuUpload className="text-gray-500" size={30} />
                                <p className="text-sm text-gray-600">Choose a file jpg, jpeg, png</p>
                            </>
                        )}
                    </label>

                    {banner && <p className="text-sm text-gray-600">{banner.name}</p>}

                    {isEditMode && !banner && eventData?.event.banner && (
                        <p className="text-xs text-gray-400">
                            Current banner will be kept if no new file is uploaded.
                        </p>
                    )}
                </div>

                <div className="flex flex-col sm:flex-row gap-4 justify-center mt-8 sm:mt-10">
                    {renderButtons()}
                </div>
            </div>

            {confirmAction && (
                <div
                    className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4"
                    onClick={() => setConfirmAction(null)}
                >
                    <div
                        className="bg-white rounded-xl p-6 max-w-sm w-full shadow-lg"
                        onClick={(e) => e.stopPropagation()}
                    >
                        <h2 className="text-lg font-bold text-gray-800 mb-2">
                            {confirmAction === "cancel" ? "Cancel Event?" : "Delete Event?"}
                        </h2>
                        <p className="text-sm text-gray-500 mb-6">
                            {confirmAction === "cancel"
                                ? "This action cannot be undone. The event will be canceled and attendees will be notified."
                                : "This will permanently delete the event draft."}
                        </p>
                        <div className="flex gap-3 justify-end">
                            <button
                                onClick={() => setConfirmAction(null)}
                                className="px-4 py-2 rounded-lg border border-gray-300 hover:bg-gray-100 text-sm"
                            >
                                Go Back
                            </button>
                            <button
                                onClick={() => {
                                    setConfirmAction(null)
                                    confirmAction === "cancel" ? handleCancel() : handleDelete()
                                }}
                                className="px-4 py-2 rounded-lg bg-red-500 hover:bg-red-600 text-white text-sm"
                            >
                                {confirmAction === "cancel" ? "Yes, Cancel" : "Yes, Delete"}
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}