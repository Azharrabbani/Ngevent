import { IoFilter } from "react-icons/io5";
import Input from "../../../../components/input";
import DateFilterDropdown from "./dateFilterDropdown";
import { SearchIcon } from "../../../../components/icon";

interface Props {
    reviewEvent: boolean;
    status: string;

    search?: string;
    setSearch?: (val: string | undefined) => void;

    sort?: string;
    setSort?: React.Dispatch<React.SetStateAction<string | undefined>>;

    dateFilter?: string;
    setDateFilter?: React.Dispatch<
        React.SetStateAction<string | undefined>
    >;

    getUpdate?: boolean;
    setGetupdate?: (
        val: boolean | undefined
    ) => void;
}

export default function EventsHeader({
    reviewEvent,
    status,
    search,
    setSearch,
    sort,
    setSort,
    dateFilter,
    setDateFilter,
    getUpdate,
    setGetupdate }: Props) {
    return (
        <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4 p-6 border-b border-gray-100">

            <div>
                <h1 className="text-xl font-semibold tracking-wide text-gray-700 uppercase">
                    {reviewEvent ? "Event Submissions" : `${status} Events`}
                </h1>
            </div>

            <div className="flex items-center gap-3 flex-wrap">

                {/* Search */}
                <div className="relative w-full lg:w-[280px]">
                    <Input
                        leftIcon={<SearchIcon />}
                        type="text"
                        placeholder="Search events, organizers..."
                        value={search}
                        onChange={(e) => setSearch?.(e.target.value)}
                        className="w-full bg-white pl-10 pr-4 py-3 rounded-lg border border-gray-300 outline-none focus:ring-2 focus:ring-blue-500 text-sm"
                    />
                </div>

                {/* Update Request Filter */}
                {reviewEvent && (
                    <button
                        onClick={() =>
                            setGetupdate?.(
                                getUpdate === true ? false : true
                            )
                        }
                        className={`px-4 py-2 rounded-lg border text-sm font-medium transition-all ${getUpdate === true
                            ? "bg-blue-50 border-blue-500 text-blue-600"
                            : "border-gray-300 text-gray-500 hover:bg-gray-50"
                            }`}
                    >
                        Update requests
                    </button>
                )}

                {/* Sort */}
                <button
                    onClick={() =>
                        setSort?.(sort === "desc" ? "asc" : "desc")
                    }
                    className="flex items-center gap-2 px-4 py-2 rounded-lg border border-gray-300 text-sm text-gray-500 hover:bg-gray-50"
                >
                    <IoFilter size={16} />

                    {sort === "desc"
                        ? "Newest first"
                        : "Oldest first"}
                </button>

                {/* Date Filter */}
                <DateFilterDropdown
                    dateFilter={dateFilter}
                    setDateFilter={setDateFilter}
                />
            </div>
        </div>
    )
}