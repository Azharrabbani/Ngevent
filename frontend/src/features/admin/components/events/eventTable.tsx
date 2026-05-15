import { IoCheckmark } from "react-icons/io5";
import { RxCross2 } from "react-icons/rx";

interface props {
    isPending: boolean;
};

export default function EventTable({ isPending }: props) {
    return (
        <div className="hidden xl:block overflow-x-auto">
            <table className="w-full min-w-[800px]">
                <thead className="bg-[#F8F9FC]">
                    <tr>
                        <th className="px-8 py-5 text-left text-sm font-semibold text-gray-600 border-b">
                            Event Details
                        </th>

                        <th className="px-8 py-5 text-left text-sm font-semibold text-gray-600 border-b">
                            Organizer
                        </th>

                        <th className="px-8 py-5 text-left text-sm font-semibold text-gray-600 border-b">
                            Submitted
                        </th>

                        <th
                            className={`
                                            px-8 py-5 text-sm font-semibold text-gray-600 border-b
                                            ${isPending ? "text-center" : "text-left"}
                                        `}
                        >
                            {isPending ? "Actions" : "Status"}
                        </th>
                    </tr>
                </thead>

                <tbody>
                    <tr className="hover:bg-gray-50 cursor-pointer transition-colors duration-200">
                        <td className="px-8 py-6 border-b border-gray-100">
                            <div className="flex items-center gap-4">
                                <div className="w-14 h-14 rounded-2xl bg-[#EEF0FF] flex items-center justify-center">
                                    🎤
                                </div>

                                <div>
                                    <h1 className="font-semibold text-lg text-gray-800">
                                        Global Tech Summit 2024
                                    </h1>

                                    <p className="text-gray-500 text-sm">
                                        Oct 15 - 17, 2024
                                    </p>
                                </div>
                            </div>
                        </td>

                        <td className="px-8 py-6 border-b border-gray-100">
                            <div>
                                <h1 className="font-semibold text-gray-800">
                                    TechCorp Inc.
                                </h1>

                                <p className="text-gray-500 text-sm">
                                    sarah@techcorp.com
                                </p>
                            </div>
                        </td>

                        <td className="px-8 py-6 border-b border-gray-100">
                            <span className="px-4 py-2 rounded-full text-sm bg-orange-100 text-orange-700">
                                2 hours ago
                            </span>
                        </td>


                        {
                            isPending ? (
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
                                    <span className="font-medium text-green-600">
                                        Active
                                    </span>
                                </td>
                            )
                        }
                    </tr>
                </tbody>
            </table>
        </div>
    )
}