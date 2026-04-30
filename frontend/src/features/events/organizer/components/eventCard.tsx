import { BsCalendar4 } from "react-icons/bs";
import { GoDotFill } from "react-icons/go";
import { MdOutlineLocationOn } from "react-icons/md";
import { cn } from "../../../../utils/cn";

type Ticket = {
  name: string;
  price: number;
  sold: number;
  total: number;
};

interface EventCardProps{
  title: string;
  date: string;
  location: string;
  status: string;
  image: string | undefined;
  revenue: number;
  tickets: Ticket[];
};

const statusColorMap  = {
  active: "text-[#0040A1]",
  pending: "text-amber-600",
  reject:  "text-red-600",
  cancel:  "text-red-600",
  done:  "text-green-600",
  draft:  "text-gray-400",
}

export default function EventCard( { 
    title,
    date, 
    location,
    status,
    image,
    revenue,
    tickets
}: EventCardProps ) {
    
    return (
        <div className="bg-white w-full max-w-sm rounded-2xl 
        shadow-sm hover:shadow-lg hover:-translate-y-1 cursor-pointer transition overflow-hidden duration-200">
            <div className="relative">
                <div className={cn(
                    "absolute flex items-center gap-1 bg-white/90 px-3 py-1 top-3 left-3 rounded-full text-sm",
                    statusColorMap[status as keyof typeof statusColorMap] || "text-gray-400"
                )}>
                    <GoDotFill/>
                    <p className="capitalize">{status}</p>
                </div>
                <img 
                    src={image === undefined ? "https://t4.ftcdn.net/jpg/16/79/44/21/360_F_1679442196_OEsi0AFKie6hYMBpvmXwwRgRYGV4U6Lz.jpg" : image}
                    alt="" 
                    className="w-full h-40 object-cover shrink-0 rounded-t-2xl"    
                />
            </div>                           

            <div className="p-5">
                <h1 className="text-xl font-bold">
                    {title}
                </h1>

                <div className="mt-2 text-[#424654] text-sm space-y-1">
                    <span className="flex items-center gap-2">
                        <BsCalendar4 className="text-[#424654] font-semibold"/>
                        <p>{date}</p>
                    </span>
                    <span className="flex items-center gap-2">
                        <MdOutlineLocationOn className="text-[#424654] font-semibold"/>
                        <p>{location}</p>
                    </span>
                </div>

                <div className="p-3 my-4 bg-[#F2F3FE] rounded-lg">
                    <p className="tex-sm text-[#424654] font-semibold">Total Revenue</p>
                    <h2 className="text-xl text-[#0040A1] font-bold">IDR {revenue.toLocaleString()}</h2>
                </div>

                <div className="flex flex-col gap-4">
                    {tickets.map((ticket, i) => {
                        const percent = Math.round((ticket.sold / ticket.total) * 100);

                        return (
                            <div key={i}>
                                <div className="flex justify-between">
                                    <h2 className="font-semibold">{ticket.name}</h2>
                                    <h2 className="text-[#0040A1] font-semibold">IDR {ticket.price.toLocaleString()}</h2>
                                </div>

                                <div className="w-full bg-gray-200 h-2 rounded-full mt-2">
                                    <div
                                        className="bg-[#0040A1] h-2 rounded-full transition-all"
                                        style={{ width: `${percent}%` }}
                                    >

                                    </div>
                                </div>

                                <div className="flex justify-between text-xs text-gray-500 mt-1">
                                    <span>{ticket.sold} / {ticket.total} Sold</span>
                                    <span>{percent}%</span>
                                </div>
                            </div>
                        );
                    }
                    )}
                </div>
            </div>    
        </div>
    );
}