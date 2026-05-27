import { CategoryIcon } from "../../../../components/icon";
import type { PaginatedData } from "../../../../types/apiResponse"
import type { categoriesPaginatedResp } from "../../../categories/types/categoryResponse"

interface Props {
    data: PaginatedData<categoriesPaginatedResp>;
};

export default function CategoriesCard({ data }: Props) {

    return (
        <div className="sm:hidden p-4 space-y-4 max-h-150 sm:max-h-125 overflow-y-auto">
            {data.rows.map((item) => {
                return (
                    <div
                        key={item.id}
                        className="hover:bg-gray-50 cursor-pointer transition-colors duration-200 border border-gray-100 rounded-2xl p-4 shadow-sm bg-white"
                    >
                        <div className="flex items-start gap-4">
                            <div className="w-14 h-14 rounded-2xl bg-[#EEF0FF] flex items-center justify-center shrink-0">
                                <CategoryIcon className="w-7 h-7 text-blue-400" />
                            </div>
                            <div className="flex-1 min-w-0">
                                <div className="flex items-start justify-between gap-3">
                                    <h1 className="font-semibold text-base text-gray-800 leading-snug">{item.name}</h1>

                                    <span className={`px-4 py-2 rounded-full text-sm bg-blue-100 font-semibold whitespace-nowrap`}>
                                        {item.total_used}
                                    </span>
                                </div>
                                <div className="mt-4">
                                    <p className="text-xs text-gray-400">Slug</p>
                                    <h2 className="text-sm font-medium text-gray-700">{item.slug}</h2>
                                </div>
                            </div>
                        </div>
                    </div>
                )
            })}

        </div>
    )
}