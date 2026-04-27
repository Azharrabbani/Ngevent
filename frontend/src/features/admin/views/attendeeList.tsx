import { useState } from "react";
import AdminSidebar from "../components/sideBar";
import { IoIosSearch } from "react-icons/io";
import { useListAttendee } from "../../profile/hooks/attendee/useListAttendee";
import { defaultPagination } from "../../../utils/pagination";
import { useAttendeeProfileDetail } from "../../profile/hooks/attendee/useAttendeeProfileDetail";
import Pagination from "../../../components/pagination";
import ProfileModal from "../components/profileModal";
import AttendeeProfile from "../components/attendeeProfile";

export default function AttendeeList() {
    const [filterAttendee, setFilterAttendee] = useState<string | undefined>();
    const [search, setSearch] = useState<string | undefined>();
    const [currentPage, setCurrentPage] = useState(1);
    const [selectedUser, setSelectedUser] = useState<string>("");
    const [isModalOpen, setIsModalOpen] = useState(false);

    const { data, isLoading } = useListAttendee({
        filter: search,
        pagination: defaultPagination(currentPage),
    });

    const { data: profile, isLoading: loadingProfile } = useAttendeeProfileDetail(selectedUser);

    const totalPage = data?.total_pages || 1;

    const handleFilter = (e: React.FormEvent) => {
        e.preventDefault();
        setSearch(filterAttendee);
        setCurrentPage(1);
    };

    return (
        <AdminSidebar>
            <>
                <h1 className="text-2xl font-semibold">Hello Admin!</h1>
                
                <form 
                    className="relative w-full max-w-md mt-5"
                    onSubmit={handleFilter}
                >
                    <input 
                        type="text" 
                        onChange={(e) => setFilterAttendee(e.target.value)}  
                        placeholder="Search..." 
                        className="w-full px-4 py-2 pr-12 rounded-lg bg-gray-200 focus:outline-none" />
                    <button
                    type="submit" 
                    className="absolute right-2 top-1/2 -translate-y-1/2 bg-blue-300 text-white p-2 rounded-full hover:bg-blue-400 transition">
                        <IoIosSearch/>
                    </button>
                </form>
                 

                <div className="mt-6 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 w-full">
                    {isLoading ? (
                        <h1>Loading....</h1>
                    ) : !data?.rows || data.rows.length === 0 ? (
                        <h1 className="text-gray-400 text-sm text-center md:text-start">
                            User not found
                        </h1>
                    ) : (
                        data?.rows.map((item) => (
                            <div
                                key={item.id}
                                onClick={() => {
                                    setSelectedUser(item.id)
                                    setIsModalOpen(true)
                                }} 
                                className="flex items-center gap-4 bg-white border-2 border-black p-5 rounded-xl 
                                            shadow-[3px_3px_0px_0px_#000] 
                                            hover:-translate-x-1 hover:-translate-y-1 hover:shadow-[6px_6px_0px_0px_#000] cursor-pointer
                                            transition-all duration-200"
                            >
                                <img 
                                    src={item.photo_profile} 
                                    alt="" 
                                    className="w-20 h-20 object-cover rounded-full shrink-0"
                                />
                                <div className="flex flex-col gap-1 overflow-hidden">
                                    <h2 className="font-bold text-lg truncate">
                                        {item.name}
                                    </h2>
    
                                    <div className="flex flex-wrap gap-2 text-sm text-gray-600">
                                        <span className="truncate">@{item.username}</span>
                                        <span className="font-medium truncate">{item.email}</span>
                                    </div>
    
                                    <div className="flex flex-wrap gap-2 text-sm text-gray-500">
                                        <span>{item.phone_number}</span>
                                        <span>{item.country}</span>
                                    </div>
                                </div>
    
                            </div>
                        )) 
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
                        setIsModalOpen(false)
                        setSelectedUser("")
                    }}
                    isLoading={loadingProfile}
                > 
                    <AttendeeProfile profile={profile}/>
                </ProfileModal>
            </>
        </AdminSidebar>
    )

}