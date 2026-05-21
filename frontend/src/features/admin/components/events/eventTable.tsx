import { IoCheckmark } from "react-icons/io5";
import { RxCross2 } from "react-icons/rx";
import type { PaginatedData } from "../../../../types/apiResponse";
import type { EventsResponse } from "../../../events/types/organizerResponse";
import { FormatRelativeTime } from "../../../../utils/formatRelativeTime";
import { getStatusColor } from "../../../../utils/status";
import { toDateString } from "../../../../utils/dateConverter";

interface Props {
    status: string;
    data: PaginatedData<EventsResponse> | undefined;
    isLoading: boolean;
    isReview: boolean;
    sort?: string
    setSort?: React.Dispatch<React.SetStateAction<string | undefined>>
}

export default function EventTable({
    status,
    isReview,
    data,
    isLoading,
}: Props) {

    return (
        <div className="hidden xl:block overflow-x-auto">
            {isLoading ? (
                <div className="p-8">
                    <h1>Loading...</h1>
                </div>
            ) : !data?.rows || data.rows.length === 0 ? (
                <div className="p-8">
                    <h1 className="text-gray-400 text-sm text-center md:text-start">
                        Events not found
                    </h1>
                </div>
            ) : (
                <table className="w-full min-w-[800px]">
                    <thead className="bg-[#F8F9FC]">
                        <tr>
                            <th className="px-8 py-5 text-left text-sm font-semibold text-gray-600 border-b">
                                Event Details
                            </th>

                            <th className="px-8 py-5 text-left text-sm font-semibold text-gray-600 border-b">
                                Organizer
                            </th>

                            <th
                                className="flex gap-2 items-center 
                                            px-8 py-5 text-left text-sm font-semibold text-gray-600 border-b"
                            >
                                Submitted
                            </th>

                            <th
                                className={`
                                    px-8 py-5 text-sm font-semibold text-gray-600 border-b
                                    ${isReview && status === "pending" ? "text-center" : "text-left"}
                                `}
                            >
                                {isReview && status === "pending" ? "Actions" : "Status"}
                            </th>
                        </tr>
                    </thead>

                    <tbody>
                        {data.rows.map((item) => {
                            const start = new Date(
                                Number(item.start_time) * 1000
                            );

                            const date = toDateString(start);

                            const submitted = FormatRelativeTime(
                                item.created_at
                            );

                            return (
                                <tr
                                    key={item.id}
                                    className="hover:bg-gray-50 cursor-pointer transition-colors duration-200"
                                >
                                    <td className="px-8 py-6 border-b border-gray-100">
                                        <div className="flex items-center gap-4">
                                            <div className="w-14 h-14 rounded-2xl bg-[#EEF0FF] overflow-hidden">
                                                <img src={item.event.banner} alt="event_banner" className="w-full h-full object-cover" />
                                            </div>

                                            <div>
                                                <h1 className="font-semibold text-lg text-gray-800">
                                                    {item.event.name}
                                                </h1>

                                                <p className="text-gray-500 text-sm">
                                                    {date}
                                                </p>
                                            </div>
                                        </div>
                                    </td>

                                    <td className="px-8 py-6 border-b border-gray-100">
                                        <div>
                                            <h1 className="font-semibold text-gray-800">
                                                {item.eo_profile.name}
                                            </h1>

                                            <p className="text-gray-500 text-sm">
                                                {item.eo_profile.email}
                                            </p>
                                        </div>
                                    </td>

                                    <td className="px-8 py-6 border-b border-gray-100">
                                        <span className="px-4 py-2 rounded-full text-sm bg-orange-100 text-orange-700">
                                            {submitted}
                                        </span>
                                    </td>

                                    {isReview && status === "pending" ? (
                                        <td className="px-8 py-6 border-b border-gray-100">
                                            <div className="flex items-center justify-center gap-3">
                                                <button className="px-7 py-3 border border-gray-300 rounded-xl hover:bg-gray-100 transition-colors">
                                                    Review
                                                </button>

                                                <button className="w-10 h-10 bg-white text-red-600 border border-red-600 rounded-lg hover:bg-red-600 hover:text-white transition-colors flex items-center justify-center">
                                                    <RxCross2 size={24} />
                                                </button>

                                                <button className="w-10 h-10 bg-blue-500 text-white rounded-lg hover:bg-blue-700 transition-colors flex items-center justify-center">
                                                    <IoCheckmark size={24} />
                                                </button>
                                            </div>
                                        </td>
                                    ) : (
                                        <td className="px-8 py-6 border-b border-gray-100">
                                            <span className={`font-medium
                                                ${getStatusColor(item?.event?.status ?? "")}`}>
                                                {item.event.status}
                                            </span>
                                        </td>
                                    )}
                                </tr>
                            );
                        })}
                    </tbody>
                </table>
            )}
        </div>
    );
}