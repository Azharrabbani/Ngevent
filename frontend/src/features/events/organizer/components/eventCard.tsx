import { GoDotFill } from "react-icons/go";
import { MdOutlineLocationOn } from "react-icons/md";
import { cn } from "../../../../utils/cn";
import { useNavigate } from "react-router-dom";
import { FiClock } from "react-icons/fi";
import { CalenderIcon } from "../../../../components/icon";
import { timeRange } from "../../../../utils/timeRange";
import { eventDateRange } from "../../../../utils/dateConverter";

interface EventCardProps {
    id: string;
    title: string;
    startDate: number;
    endDate: number;
    startTime: number;
    endTime: number;
    location: string;
    status: string;
    image: string | undefined;
}

const statusColorMap = {
    active: "text-[#0040A1]",
    pending: "text-amber-600",
    rejected: "text-red-600",
    cancelled: "text-red-600",
    done: "text-green-600",
    draft: "text-gray-400",
};

export default function EventCard({
    id,
    title,
    startDate,
    endDate,
    startTime,
    endTime,
    location,
    status,
    image,
}: EventCardProps) {
    const navigate = useNavigate();

    const start = new Date(Number(startTime) * 1000);
    const end = new Date(Number(endTime) * 1000);

    const eventDate = eventDateRange(
        startDate,
        endDate
    );

    const timeRangeVal = timeRange(start, end);

    return (
        <div
            onClick={() => navigate(`/organizer/event/view/${id}`)}
            className="bg-white w-full max-w-sm rounded-2xl
        shadow-sm hover:shadow-lg hover:-translate-y-1 cursor-pointer transition overflow-hidden duration-200"
        >
            <div className="relative">
                <div
                    className={cn(
                        "absolute flex items-center gap-1 bg-white/90 px-3 py-1 top-3 left-3 rounded-full text-sm",
                        statusColorMap[status as keyof typeof statusColorMap] || "text-gray-400"
                    )}
                >
                    <GoDotFill />
                    <p className="capitalize">{status}</p>
                </div>
                <img
                    src={
                        image === undefined
                            ? "https://t4.ftcdn.net/jpg/16/79/44/21/360_F_1679442196_OEsi0AFKie6hYMBpvmXwwRgRYGV4U6Lz.jpg"
                            : image
                    }
                    alt=""
                    className="w-full h-40 object-cover shrink-0 rounded-t-2xl"
                />
            </div>

            <div className="p-5">
                <h1 className="text-xl font-bold">{title}</h1>

                <div className="mt-2 text-[#424654] text-sm space-y-1">
                    <span className="flex items-start gap-2">
                        <div className="flex flex-col">
                            <div className="flex items-center gap-2">
                                <CalenderIcon className="text-[#424654] font-semibold mt-0.5 shrink-0" />
                                <p>{eventDate}</p>
                            </div>
                            <div className="flex items-center gap-2">
                                <FiClock className="text-[#424654] font-semibold mt-0.5 shrink-0" />
                                <p>{timeRangeVal} WIB</p>
                            </div>
                        </div>
                    </span>
                    <span className="flex items-center gap-2">
                        <MdOutlineLocationOn className="text-[#424654] font-semibold" />
                        <p>{location}</p>
                    </span>
                </div>
            </div>
        </div>
    );
}