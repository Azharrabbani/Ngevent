import { HiOutlineDownload } from "react-icons/hi";

interface Props {
    previewUrl: string;
    periodLabel: string;
    filename: string;
    onDownload: () => void;
}

export default function PreviewPanel({ previewUrl, periodLabel, filename, onDownload }: Props) {
    return (
        <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
            <div className="flex items-center justify-between px-6 py-4 border-b border-gray-100">
                <div>
                    <p className="text-sm font-semibold text-[#1e293b]">
                        Preview — {periodLabel}
                    </p>
                    <p className="text-xs text-gray-400 mt-0.5">{filename}</p>
                </div>
                <button
                    onClick={onDownload}
                    className="flex items-center gap-2 bg-[#0056D2] hover:bg-[#0046b0] text-white text-sm font-semibold px-5 py-2.5 rounded-xl transition-all duration-200 shadow-sm"
                >
                    <HiOutlineDownload className="text-base" />
                    Download
                </button>
            </div>

            <iframe
                src={previewUrl}
                title="Report Preview"
                className="w-full"
                style={{ height: "78vh", border: "none" }}
            />
        </div>
    );
}