import { FiChevronLeft, FiChevronRight } from "react-icons/fi";

interface Props {
    currentPage: number;
    totalPages: number;
    onPageChange: (page: number) => void;
}

export default function PaginationTabs({ currentPage, totalPages, onPageChange }: Props) {
    if (totalPages <= 1) return null;

    const getPages = (): (number | "...")[] => {
        if (totalPages <= 7) {
            return Array.from({ length: totalPages }, (_, i) => i + 1);
        }

        const pages: (number | "...")[] = [1];

        if (currentPage > 3) pages.push("...");

        const start = Math.max(2, currentPage - 1);
        const end = Math.min(totalPages - 1, currentPage + 1);
        for (let i = start; i <= end; i++) pages.push(i);

        if (currentPage < totalPages - 2) pages.push("...");
        pages.push(totalPages);

        return pages;
    };

    const pages = getPages();

    return (
        <div className="flex items-center justify-center gap-1.5 mt-10">
            <button
                onClick={() => onPageChange(currentPage - 1)}
                disabled={currentPage === 1}
                className="
                    w-9 h-9 flex items-center justify-center
                    rounded-lg border border-slate-200 bg-white
                    text-slate-500 hover:text-blue-600 hover:border-blue-300
                    disabled:opacity-40 disabled:cursor-not-allowed
                    transition
                "
            >
                <FiChevronLeft className="w-4 h-4" />
            </button>

            {pages.map((page, i) =>
                page === "..." ? (
                    <span
                        key={`ellipsis-${i}`}
                        className="w-9 h-9 flex items-center justify-center text-slate-400 text-sm"
                    >
                        …
                    </span>
                ) : (
                    <button
                        key={page}
                        onClick={() => onPageChange(page as number)}
                        className={`
                            w-9 h-9 flex items-center justify-center rounded-lg text-sm font-medium
                            border transition
                            ${currentPage === page
                                ? "bg-slate-800 border-slate-800 text-white shadow-sm"
                                : "bg-white border-slate-200 text-slate-600 hover:text-blue-600 hover:border-blue-300"
                            }
                        `}
                    >
                        {page}
                    </button>
                )
            )}

            <button
                onClick={() => onPageChange(currentPage + 1)}
                disabled={currentPage === totalPages}
                className="
                    w-9 h-9 flex items-center justify-center
                    rounded-lg border border-slate-200 bg-white
                    text-slate-500 hover:text-blue-600 hover:border-blue-300
                    disabled:opacity-40 disabled:cursor-not-allowed
                    transition
                "
            >
                <FiChevronRight className="w-4 h-4" />
            </button>
        </div>
    );
}