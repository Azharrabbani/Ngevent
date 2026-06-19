import { HiOutlineDocumentReport } from "react-icons/hi";

export default function ReportHeader() {
    return (
        <div className="mb-8">
            <div className="flex items-center gap-3 mb-1">
                <HiOutlineDocumentReport className="text-[#0056D2] text-3xl" />
                <h1 className="text-2xl font-bold text-[#1e293b]">Event Report</h1>
            </div>
            <p className="text-sm text-gray-500 ml-10">
                Generate and download your event report
            </p>
        </div>
    );
}