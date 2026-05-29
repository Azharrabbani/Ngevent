import { useState } from "react";
import AdminSidebar from "../components/sideBar";
import { IoIosSearch } from "react-icons/io";
import { defaultPagination } from "../../../utils/pagination";
import Pagination from "../../../components/pagination";
import Button from "../../../components/Button";
import { PlusIcon, SpinnerIcon, UserVerifiedIcon } from "../../../components/icon";
import CreateAdminSlideOver from "../components/createAdminSlideOver";
import { toDateString } from "../../../utils/dateConverter";
import { useAuth } from "../../../lib/auth";
import { useListUsers } from "../../profile/hooks/useListUsers";

export default function AdminList() {
    const [filterAdmin, setFilterAdmin] = useState<string | undefined>();
    const [search, setSearch] = useState<string | undefined>();
    const [currentPage, setCurrentPage] = useState(1);
    const [isSlideOverOpen, setIsSlideOverOpen] = useState(false);

    const { data, isLoading } = useListUsers({
        role: "admin",
        email: search,
        pagination: defaultPagination(currentPage),
    })

    const totalPage = data?.total_pages || 1;

    const handleFilter = (e: React.FormEvent) => {
        e.preventDefault();
        setSearch(filterAdmin);
        setCurrentPage(1);
    };

    const { user, loading } = useAuth()

    if (loading) {
        return null
    }

    const sortedUsers = [...(data?.rows || [])].sort((a, b) => {
        if (a.id === user?.id) return -1;
        if (b.id === user?.id) return 1;
        return 0;
    });

    return (
        <AdminSidebar>
            <>
                <h1 className="text-2xl font-semibold">Hello Admin!</h1>

                <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
                    <form
                        className="relative w-full max-w-md mt-5"
                        onSubmit={handleFilter}
                    >
                        <input
                            type="text"
                            onChange={(e) => setFilterAdmin(e.target.value)}
                            placeholder="Search..."
                            className="w-full px-4 py-2 pr-12 rounded-lg bg-gray-200 focus:outline-none" />
                        <button
                            type="submit"
                            className="absolute right-2 top-1/2 -translate-y-1/2 bg-blue-300 text-white p-2 rounded-full hover:bg-blue-400 transition">
                            <IoIosSearch />
                        </button>
                    </form>

                    <Button
                        className="px-6 py-3 rounded-xl bg-[#0066FF]"
                        onClick={() => setIsSlideOverOpen(true)}
                    >
                        <div className="flex justify-center items-center gap-2">
                            <PlusIcon />
                            <p>Register new admin</p>
                        </div>
                    </Button>

                    <CreateAdminSlideOver
                        isOpen={isSlideOverOpen}
                        onClose={() => setIsSlideOverOpen(false)}
                    />

                </div>


                <div className="mt-6 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 w-full">
                    {isLoading ? (
                        <div className="flex justify-center py-20">
                            <SpinnerIcon className="animate-spin w-8 h-8 text-blue-500" />
                        </div>
                    ) : !data?.rows || data.rows.length === 0 ? (
                        <div className="flex flex-col items-center justify-center py-20 text-center">
                            <h1 className="text-xl font-semibold text-gray-700">
                                No admin found
                            </h1>
                        </div>
                    ) : (
                        sortedUsers.map((item) => (
                            <div
                                key={item.id}
                                className="
                                    flex flex-col sm:flex-row
                                    items-start sm:items-center
                                    gap-4
                                    bg-white border-2 border-black
                                    p-5 rounded-xl
                                    shadow-[3px_3px_0px_0px_#000]
                                    hover:-translate-x-1 hover:-translate-y-1
                                    hover:shadow-[6px_6px_0px_0px_#000]
                                    transition-all duration-200
                                "
                            >
                                <UserVerifiedIcon
                                    className="w-16 h-16 sm:w-20 sm:h-20 shrink-0"
                                />
                                <div className="flex flex-col gap-1 overflow-hidden">
                                    <div className="flex items-center gap-2 flex-wrap">
                                        <h2 className="font-bold text-lg break-all">
                                            {item.email}
                                        </h2>

                                        {user?.id === item.id && (
                                            <span
                                                className="
                                                    bg-blue-500 text-white text-[10px] font-semibold
                                                    px-2 py-1 rounded-full whitespace-nowrap
                                                "
                                            >
                                                You
                                            </span>
                                        )}
                                    </div>

                                    <div className="flex flex-wrap gap-2 text-sm text-gray-600">
                                        <span className="truncate">id: {item.id}</span>
                                    </div>

                                    <div className="flex flex-wrap gap-2 text-sm text-gray-500">
                                        <span>Register at: {toDateString(new Date(Number(item?.created_at) * 1000), "long")}</span>
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
            </>
        </AdminSidebar>
    )

}