import { IoIosSearch } from "react-icons/io";
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
import { IoIosArrowDown } from "react-icons/io";

export default function OrganizerList() {
    const [filterOrganizer, setFilterOrganizer] = useState<string | undefined>();
    const [search, setSearch] = useState<string | undefined>();
    const [currentPage, setCurrentPage] = useState(1);
    const [status, setStatus] = useState<string | null>(null);
    const [openStatus, setOpenStatus] = useState(false);
    const [selectedUser, setSelectedUser] = useState<string>("");
    const [isModalOpen, setIsModalOpen] = useState(false);

    const dropdownRef = useRef<HTMLDivElement>(null);

    const { data, isLoading } = useListOrganizer({
        filter: search,
        status: status,
        pagination: defaultPagination(currentPage),
    });

    const { data: profile, isLoading: loadingProfile } = useOrganizerProfileDetail(selectedUser);

    const { data: organizerUpdate, isError, isPending: loadingUpdate } = useGetOrganizerUpdate(selectedUser);

    const totalPage = data?.total_pages || 1;

    useEffect(() => {
        function handleClickOutside(e: MouseEvent) {
            if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) {
                setOpenStatus(false);
            }
        }
        document.addEventListener("mousedown", handleClickOutside);
        return () => document.removeEventListener("mousedown", handleClickOutside);
    }, [search, status]);

    const handleFilter = (e: React.FormEvent) => {
        e.preventDefault();

        setSearch(
            filterOrganizer?.trim()
                ? filterOrganizer.trim()
                : undefined
        );

        setCurrentPage(1);
    };

    return (
        <AdminSidebar>
            <>
                <h1 className="text-2xl font-semibold">Hello Admin!</h1>

                <div className="flex flex-col md:flex-row gap-4 w-full mt-5">
                    <form
                        className="relative w-full max-w-md"
                        onSubmit={handleFilter}
                    >
                        <input
                            type="text"
                            onChange={(e) => setFilterOrganizer(e.target.value)}
                            placeholder="Search..."
                            className="w-full px-4 py-2 pr-12 rounded-lg bg-gray-200 focus:outline-none" />
                        <button
                            type="submit"
                            className="absolute right-2 top-1/2 -translate-y-1/2 bg-blue-400 text-white p-2 rounded-full hover:bg-blue-500 transition">
                            <IoIosSearch />
                        </button>

                    </form>

                    <div className="relative w-full md:w-auto" ref={dropdownRef}>
                        <button
                            type="button"
                            onClick={() => setOpenStatus(!openStatus)}
                            className="w-full md:w-auto border border-gray-300 bg-white text-gray-600 px-4 py-2 rounded-lg flex justify-between items-center gap-2"
                        >
                            {status ? status : "Status"}
                            <IoIosArrowDown className={`transition ${openStatus ? "rotate-180" : ""}`} />
                        </button>

                        {openStatus && (
                            <div
                                className="absolute z-10 mt-2 w-full md:w-40 bg-white border border-gray-200 rounded-lg shadow-md overflow-hidden"
                            >
                                {["approved", "pending", "rejected", "deactivated"].map((item) => (
                                    <button
                                        key={item}
                                        type="button"
                                        onClick={() => {
                                            setStatus(item);
                                            setOpenStatus(false);
                                        }}
                                        className={`w-full text-left px-4 py-2 hover:bg-gray-100 ${status === item ? "bg-gray-100 font-semibold" : ""}`}
                                    >
                                        {item.charAt(0).toUpperCase() + item.slice(1)}
                                    </button>
                                ))}

                                <button
                                    type="button"
                                    onClick={() => {
                                        setStatus("approved");
                                        setOpenStatus(false);
                                    }}
                                    className="w-full text-left px-4 py-2 text-red-500 hover:bg-gray-100"
                                >
                                    Clear
                                </button>
                            </div>
                        )}
                    </div>

                </div>

                <div className="mt-6 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 w-full">
                    {isLoading ? (
                        <h1>Loading....</h1>
                    ) : !data?.rows || data.rows.length === 0 ? (
                        <h1 className="text-gray-400 text-sm text-center md:text-start">
                            User not found
                        </h1>
                    ) :
                        (
                            data.rows.map((item) => (
                                <div
                                    key={item.id}
                                    onClick={() => {
                                        setSelectedUser(item.id)
                                        setIsModalOpen(true)
                                    }}
                                    className="group relative flex items-center gap-4 bg-white border-2 border-black p-5 rounded-xl 
                                                shadow-[3px_3px_0px_0px_#000] 
                                                hover:-translate-x-1 hover:-translate-y-1 hover:shadow-[6px_6px_0px_0px_#000] cursor-pointer
                                                transition-all duration-200"
                                >
                                    <div className={`
                                        absolute left-0 top-0 h-full 
                                        w-2 group-hover:w-28 
                                    ${item.status.status === "approved" ?
                                            "bg-blue-500 " :
                                            item.status.status === "pending" ? "bg-amber-500" : "bg-red-500"}
                                        rounded-l-xl
                                        transition-all duration-300
                                        flex items-center justify-center
                                        overflow-hidden
                                    `}
                                    >
                                        <span className="
                                            text-white text-sm font-semibold
                                            opacity-0
                                            group-hover:opacity-100
                                            transition-opacity duration-200
                                            whitespace-nowrap
                                        ">
                                            {item.status.status}
                                        </span>
                                    </div>
                                    <img
                                        src={item.photo_profile}
                                        alt=""
                                        className="w-20 h-20 object-cover rounded-full shrink-0"
                                    />

                                    <div className="flex flex-col gap-1 overflow-hidden">
                                        <div className="flex items-center gap-1">
                                            <h2 className="font-bold text-lg truncate">
                                                {item.name}
                                            </h2>
                                            {item.status.status === "approved" && <MdVerified className="text-blue-600" />}
                                        </div>

                                        <p className={`
                                            text-xs md:hidden
                                            ${item.status.status === "approved" ?
                                                "text-blue-500" :
                                                item.status.status === "pending" ? "text-amber-500" :
                                                    "text-red-500"} 
                                        `}>
                                            {item.status.status}
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
                            ))
                        )
                    }

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
                        setIsModalOpen(false)
                        setSelectedUser("")
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
                            setIsModalOpen(false)
                            setSelectedUser("")
                        }}
                    />

                </ProfileModal>

            </>
        </AdminSidebar>
    )
}