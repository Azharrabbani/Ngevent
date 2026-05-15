import { IoCheckmark } from "react-icons/io5";
import { RxCross2 } from "react-icons/rx";

interface Props {
    isPending: boolean;
};

export default function EventCard({ isPending }: Props) {
    return (
        <div className="xl:hidden p-4 space-y-4">
            <div className="border border-gray-100 rounded-2xl p-4 shadow-sm bg-white">
                <div className="flex items-start gap-4">
                    <div className="w-14 h-14 rounded-2xl bg-[#EEF0FF] flex items-center justify-center shrink-0">
                        🎤
                    </div>

                    <div className="flex-1 min-w-0">
                        <div className="flex items-start justify-between gap-3">
                            <h1 className="font-semibold text-base text-gray-800 leading-snug">
                                Global Tech Summit 2024
                            </h1>

                            {!isPending && (
                                <span className="text-sm font-medium text-green-600 whitespace-nowrap">
                                    Active
                                </span>
                            )}
                        </div>

                        <p className="text-sm text-gray-500 mt-1">
                            Oct 15 - 17, 2024
                        </p>

                        <div className="mt-4">
                            <p className="text-xs text-gray-400">
                                Organizer
                            </p>

                            <h2 className="text-sm font-medium text-gray-700">
                                TechCorp Inc.
                            </h2>

                            <p className="text-sm text-gray-500 break-all">
                                sarah@techcorp.com
                            </p>
                        </div>

                        <div className="mt-5 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
                            <span className="w-fit px-3 py-1 rounded-full text-xs bg-orange-100 text-orange-700">
                                2 hours ago
                            </span>

                            {isPending && (
                                <div className="flex items-center gap-2">
                                    <button className="flex-1 sm:flex-none px-5 py-2 border border-gray-300 rounded-xl text-sm hover:bg-gray-100 transition-colors">
                                        Review
                                    </button>

                                    <button className="w-10 h-10 bg-white text-red-600 border border-red-600 rounded-lg hover:bg-red-600 hover:text-white transition-colors flex items-center justify-center">
                                        <RxCross2 size={20} />
                                    </button>

                                    <button className="w-10 h-10 bg-blue-500 text-white rounded-lg hover:bg-blue-700 transition-colors flex items-center justify-center">
                                        <IoCheckmark size={20} />
                                    </button>
                                </div>
                            )}
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
}