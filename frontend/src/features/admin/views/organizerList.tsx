import { IoIosSearch, IoIosArrowDown } from "react-icons/io";
import AdminSidebar from "../components/sideBar";
import { MdVerified } from "react-icons/md";
import { defaultPagination } from "../../../utils/pagination";
import { useListOrganizer } from "../../profile/hooks/organizer/useListOrganizer";
import { useEffect, useRef, useState } from "react";
import Pagination from "../../../components/pagination";
import { useOrganizerProfileDetail } from "../../profile/hooks/organizer/useOrganizerProfileDetail";
import ProfileModal from "../components/profiles/profileModal";
import OrganizerProfile from "../components/profiles/organizerProfile";
import { useGetOrganizerUpdate } from "../../profile/hooks/organizer/useGetOrganizerUpdate";

type FilterMode = "status" | "updates";

interface ActiveFilter {
    mode: FilterMode;
    value: string;
}

const STATUS_OPTIONS = ["approved", "pending", "rejected", "deactivated"] as const;

function accentClass(status: string, hasUpdates: boolean): string {
    if (hasUpdates) return "bg-teal-500";
    switch (status) {
        case "approved": return "bg-blue-500";
        case "pending": return "bg-amber-500";
        case "rejected": return "bg-red-500";
        case "deactivated": return "bg-red-500";
        default: return "bg-gray-400";
    }
}

function statusTextClass(status: string, hasUpdates: boolean): string {
    if (hasUpdates) return "text-teal-500";
    switch (status) {
        case "approved": return "text-blue-500";
        case "pending": return "text-amber-500";
        case "rejected":
        case "deactivated": return "text-red-500";
        default: return "text-gray-400";
    }
}

function statusLabel(status: string, hasUpdates: boolean): string {
    if (hasUpdates) return "Has Updates";
    return status.charAt(0).toUpperCase() + status.slice(1);
}

export default function OrganizerList() {
    const [filterOrganizer, setFilterOrganizer] = useState<string | undefined>();
    const [search, setSearch] = useState<string | undefined>();
    const [currentPage, setCurrentPage] = useState(1);
    const [activeFilter, setActiveFilter] = useState<ActiveFilter | null>(null);
    const [openStatus, setOpenStatus] = useState(false);
    const [selectedUser, setSelectedUser] = useState<string>("");
    const [isModalOpen, setIsModalOpen] = useState(false);

    const dropdownRef = useRef<HTMLDivElement>(null);

    const apiStatus = activeFilter?.mode === "status" ? activeFilter.value : null;
    const apiRequestUpdates = activeFilter?.mode === "updates" ? true : undefined;

    const { data, isLoading } = useListOrganizer({
        filter: search,
        status: apiStatus,
        request_updates: apiRequestUpdates,
        pagination: defaultPagination(currentPage),
    });

    const { data: profile, isLoading: loadingProfile } =
        useOrganizerProfileDetail(selectedUser);

    const { data: organizerUpdate, isError, isPending: loadingUpdate } =
        useGetOrganizerUpdate(selectedUser);

    const totalPage = data?.total_pages || 1;

    useEffect(() => {
        function handleClickOutside(e: MouseEvent) {
            if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) {
                setOpenStatus(false);
            }
        }
        document.addEventListener("mousedown", handleClickOutside);
        return () => document.removeEventListener("mousedown", handleClickOutside);
    }, []);

    const handleFilter = (e: React.FormEvent) => {
        e.preventDefault();
        setSearch(filterOrganizer?.trim() || undefined);
        setCurrentPage(1);
    };

    const handleSelectFilter = (mode: FilterMode, value: string) => {
        setActiveFilter({ mode, value });
        setCurrentPage(1);
        setOpenStatus(false);
    };

    const handleClearFilter = () => {
        setActiveFilter(null);
        setCurrentPage(1);
        setOpenStatus(false);
    };

    const dropdownLabel = activeFilter
        ? activeFilter.value === "updates"
            ? "Has Updates"
            : activeFilter.value.charAt(0).toUpperCase() + activeFilter.value.slice(1)
        : "Filter";

    return (
        <AdminSidebar>
            <div className="p-6 w-full">
                <h1 className="text-2xl font-semibold">Event Owner Management</h1>

                <div className="flex flex-col md:flex-row gap-4 w-full mt-5">
                    <form className="relative w-full max-w-md" onSubmit={handleFilter}>
                        <input
                            type="text"
                            onChange={(e) => setFilterOrganizer(e.target.value)}
                            placeholder="Search..."
                            className="w-full px-4 py-2 pr-12 rounded-lg bg-gray-200 focus:outline-none"
                        />
                        <button
                            type="submit"
                            className="absolute right-2 top-1/2 -translate-y-1/2 bg-blue-400 text-white p-2 rounded-full hover:bg-blue-500 transition"
                        >
                            <IoIosSearch />
                        </button>
                    </form>

                    {/* Filter dropdown */}
                    <div className="relative w-full md:w-auto" ref={dropdownRef}>
                        <button
                            type="button"
                            onClick={() => setOpenStatus(!openStatus)}
                            className={`
                                w-full md:w-auto border px-4 py-2 rounded-lg flex justify-between items-center gap-2 transition-colors
                                ${activeFilter
                                    ? activeFilter.value === "updates"
                                        ? "border-teal-400 bg-teal-50 text-teal-700"
                                        : "border-blue-300 bg-blue-50 text-blue-700"
                                    : "border-gray-300 bg-white text-gray-600"}
                            `}
                        >
                            {dropdownLabel}
                            <IoIosArrowDown className={`transition ${openStatus ? "rotate-180" : ""}`} />
                        </button>

                        {openStatus && (
                            <div className="absolute z-10 mt-2 w-full md:w-44 bg-white border border-gray-200 rounded-lg shadow-md overflow-hidden">
                                {/* Status options */}
                                {STATUS_OPTIONS.map((item) => (
                                    <button
                                        key={item}
                                        type="button"
                                        onClick={() => handleSelectFilter("status", item)}
                                        className={`w-full text-left px-4 py-2 hover:bg-gray-100 transition-colors ${activeFilter?.mode === "status" && activeFilter.value === item
                                                ? "bg-gray-100 font-semibold"
                                                : ""
                                            }`}
                                    >
                                        {item.charAt(0).toUpperCase() + item.slice(1)}
                                    </button>
                                ))}

                                <div className="border-t border-gray-100 my-1" />

                                {/* Updates filter */}
                                <button
                                    type="button"
                                    onClick={() => handleSelectFilter("updates", "updates")}
                                    className={`w-full text-left px-4 py-2 flex items-center gap-2 hover:bg-teal-50 transition-colors ${activeFilter?.mode === "updates"
                                            ? "bg-teal-50 font-semibold text-teal-700"
                                            : "text-teal-600"
                                        }`}
                                >
                                    <span className="w-2 h-2 rounded-full bg-teal-500 shrink-0" />
                                    Has Updates
                                </button>

                                {activeFilter && (
                                    <>
                                        <div className="border-t border-gray-100 my-1" />
                                        <button
                                            type="button"
                                            onClick={handleClearFilter}
                                            className="w-full text-left px-4 py-2 text-red-500 hover:bg-gray-100 transition-colors"
                                        >
                                            Clear filter
                                        </button>
                                    </>
                                )}
                            </div>
                        )}
                    </div>
                </div>

                <div className="flex flex-wrap gap-3 mt-4 text-xs text-gray-500">
                    {[
                        { color: "bg-blue-500", label: "Approved" },
                        { color: "bg-amber-500", label: "Pending" },
                        { color: "bg-red-500", label: "Rejected / Deactivated" },
                        { color: "bg-teal-500", label: "Has Updates" },
                    ].map(({ color, label }) => (
                        <span key={label} className="flex items-center gap-1.5">
                            <span className={`w-2.5 h-2.5 rounded-full ${color}`} />
                            {label}
                        </span>
                    ))}
                </div>

                {/* Cards */}
                <div className="mt-6 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 w-full">
                    {isLoading ? (
                        <p className="text-gray-400 text-sm">Loading...</p>
                    ) : !data?.rows || data.rows.length === 0 ? (
                        <p className="text-gray-400 text-sm text-center md:text-start">
                            No event owners found
                        </p>
                    ) : (
                        data.rows.map((item) => {
                            const hasUpdates = item.request_updates === true;

                            return (
                                <div
                                    key={item.id}
                                    onClick={() => {
                                        setSelectedUser(item.id);
                                        setIsModalOpen(true);
                                    }}
                                    className="group relative flex items-center gap-4 bg-white border-2 border-black p-5 rounded-xl
                                               shadow-[3px_3px_0px_0px_#000]
                                               hover:-translate-x-1 hover:-translate-y-1 hover:shadow-[6px_6px_0px_0px_#000] cursor-pointer
                                               transition-all duration-200"
                                >
                                    <div
                                        className={`
                                            absolute left-0 top-0 h-full
                                            w-2 group-hover:w-28
                                            ${accentClass(item.status.status, hasUpdates)}
                                            rounded-l-xl
                                            transition-all duration-300
                                            flex items-center justify-center
                                            overflow-hidden
                                        `}
                                    >
                                        <span className="text-white text-sm font-semibold opacity-0 group-hover:opacity-100 transition-opacity duration-200 whitespace-nowrap">
                                            {statusLabel(item.status.status, hasUpdates)}
                                        </span>
                                    </div>

                                    <img
                                        src={item.photo_profile}
                                        alt=""
                                        className="w-20 h-20 object-cover rounded-full shrink-0"
                                    />

                                    <div className="flex flex-col gap-1 overflow-hidden">
                                        <div className="flex items-center gap-1">
                                            <h2 className="font-bold text-lg truncate">{item.name}</h2>
                                            {item.status.status === "approved" && !hasUpdates && (
                                                <MdVerified className="text-blue-600" />
                                            )}
                                            {hasUpdates && (
                                                <MdVerified className="text-teal-500" />
                                            )}
                                        </div>

                                        {/* Mobile status label */}
                                        <p className={`text-xs md:hidden ${statusTextClass(item.status.status, hasUpdates)}`}>
                                            {statusLabel(item.status.status, hasUpdates)}
                                        </p>

                                        <div className="flex flex-wrap gap-2 text-sm text-gray-600">
                                            <span className="font-medium truncate">{item.email}</span>
                                        </div>

                                        <div className="flex flex-wrap gap-2 text-sm text-gray-500">
                                            <span>{item.phone_number}</span>
                                            <span>{item.country}</span>
                                        </div>

                                        <div className="flex flex-wrap gap-2 text-sm text-gray-500">
                                            <span>npwp: {item.company_detail.npwp}</span>
                                            <span>nib: {item.company_detail.nib}</span>
                                        </div>
                                    </div>
                                </div>
                            );
                        })
                    )}
                </div>

                <Pagination
                    currentPage={currentPage}
                    totalPage={totalPage}
                    onPrev={() => setCurrentPage((prev) => Math.max(prev - 1, 1))}
                    onNext={() => setCurrentPage((prev) => Math.min(prev + 1, totalPage))}
                    onCurrent={(page) => setCurrentPage(page)}
                />

                <ProfileModal
                    isOpen={isModalOpen}
                    onClose={() => {
                        setIsModalOpen(false);
                        setSelectedUser("");
                    }}
                    isLoading={loadingProfile}
                >
                    <OrganizerProfile
                        profile={profile}
                        update={organizerUpdate}
                        profileLoading={loadingProfile}
                        updateLoading={loadingUpdate}
                        isError={isError}
                        onClose={() => {
                            setIsModalOpen(false);
                            setSelectedUser("");
                        }}
                    />
                </ProfileModal>
            </div>
        </AdminSidebar>
    );
}