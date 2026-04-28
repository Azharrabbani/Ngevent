import { getPage } from "../utils/pagination";

interface Props {
    onPrev: () => void;
    onCurrent: (page: number) => void;
    onNext: () => void;
    currentPage: number;
    totalPage: number;

}

export default function Pagination({
    onPrev,
    onCurrent,
    onNext,
    currentPage,
    totalPage
}: Props) {

    const pages = getPage(currentPage, totalPage);

    return (
        <div className="flex justify-center my-8">
            <div className="flex items-center gap-1 sm:gap-2 bg-white/70 backdrop-blur px-2 py-2 rounded-xl shadow-md">
                <button
                    onClick={onPrev}
                    className="px-3 py-1 rounded-lg hover:bg-gray-200 disabled:opacity-50"
                    disabled={currentPage === 1}
                >
                    {"<"}
                </button>

                {pages.map((page) => (
                    <button
                        key={page}
                        onClick={() => onCurrent(page)}
                        className={`px-3 py-1 rounded-lg text-sm sm:text-base transition
                        ${currentPage === page
                            ? "bg-blue-500 text-white shadow"
                            : "hover:bg-gray-200"
                        }`}
                    >
                        {page}
                    </button>
                ))}

                <button
                    onClick={onNext}
                    className="px-3 py-1 rounded-lg hover:bg-gray-200 disabled:opacity-50"
                    disabled={currentPage === totalPage}
                >
                    {">"}
                </button>

            </div>
        </div>
    );
}