import type { PaginatedData } from "../../../../types/apiResponse";
import type { categoriesPaginatedResp } from "../../../categories/types/categoryResponse";
import CategoriesTable from "./categoriesTable";
import CategoriesCard from "./categoriesCard";
import Pagination from "../../../../components/pagination";
import { SpinnerIcon } from "../../../../components/icon";

interface Props {
    data: PaginatedData<categoriesPaginatedResp> | undefined;
    isLoading: boolean;
    search?: string;
    setSearch?: (val: string | undefined) => void;
    currentPage: number;
    totalPage: number;
    setCurrentPage: React.Dispatch<React.SetStateAction<number>>;
};


export default function CategoriesList({ data, isLoading, currentPage, totalPage, setCurrentPage }: Props) {
    if (isLoading) {
        return (
            <div className="flex justify-center py-20">
                <SpinnerIcon className="animate-spin w-8 h-8 text-blue-500" />
            </div>
        );
    }

    const isEmpty = !data || data.total_rows === 0;

    if (isEmpty) {
        return (
            <div className="flex flex-col items-center justify-center py-20 text-center">
                <h1 className="text-xl font-semibold text-gray-700">
                    No categories found
                </h1>
            </div>
        );
    }

    return (
        <div>
            <CategoriesTable data={data} />

            <CategoriesCard data={data} />

            {data && (
                <div className="flex flex-col md:flex-row items-center justify-between gap-4 px-6 md:px-8 border-t border-gray-100">

                    <p className="text-sm text-gray-500">
                        Showing 1 to 4 of {data?.total_rows} categories
                    </p>

                    <Pagination
                        currentPage={currentPage}
                        totalPage={totalPage}
                        onPrev={() =>
                            setCurrentPage((prev) =>
                                Math.max(prev - 1, 1)
                            )
                        }
                        onNext={() =>
                            setCurrentPage((prev) =>
                                Math.min(prev + 1, totalPage)
                            )
                        }
                        onCurrent={(page) => setCurrentPage(page)}
                    />
                </div>
            )}
        </div>
    )
}