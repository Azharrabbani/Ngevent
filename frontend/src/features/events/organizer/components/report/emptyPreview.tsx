import { HiOutlineDocumentReport } from "react-icons/hi";

export default function EmptyPreview() {
    return (
        <div className="bg-white rounded-2xl border border-dashed border-gray-200 flex flex-col items-center justify-center py-24 text-center">
            <div className="bg-[#0056D2]/8 rounded-full p-5 mb-5">
                <HiOutlineDocumentReport className="text-[#0056D2] text-5xl" />
            </div>
            <p className="text-[#1e293b] font-semibold text-lg mb-1">
                No preview yet
            </p>
            <p className="text-gray-400 text-sm max-w-xs">
                Select your period settings above and click{" "}
                <span className="font-semibold text-[#0056D2]">Generate Report</span>{" "}
                to preview and download your PDF.
            </p>
        </div>
    );
}