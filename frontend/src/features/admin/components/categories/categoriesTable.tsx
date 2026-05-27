
import { CategoryIcon } from "../../../../components/icon"
import type { PaginatedData } from "../../../../types/apiResponse"
import type { categoriesPaginatedResp } from "../../../categories/types/categoryResponse"

interface Props {
    data: PaginatedData<categoriesPaginatedResp>;
}

export default function CategoriesTable({ data }: Props) {

    return (
        <div className="hidden sm:flex justify-center w-full">
            <div className="w-full lg:max-w-8xl overflow-hidden border-2 border-[#C2C6D8] rounded-xl">
                <table className="w-full table-auto">
                    <thead className="bg-[#EFF4FF]">
                        <tr>
                            <th className="px-8 py-5 text-left text-md font-semibold text-gray-600 border-b-2 border-[#C2C6D8]">CATEGORY NAME</th>
                            <th className="px-8 py-5 text-left text-md font-semibold text-gray-600 border-b-2 border-[#C2C6D8]">SLUG</th>
                            <th className="w-[180px] px-8 py-5 text-left text-md font-semibold text-gray-600 border-b-2 border-[#C2C6D8]">
                                TOTAL USED
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        {data.rows.map((item) => {
                            return (
                                <tr
                                    key={item.id}
                                    className="hover:bg-gray-100 cursor-pointer transition-colors duration-200 border-b-2 border-[#C2C6D8] last:border-0"
                                >
                                    <td className="px-8 py-6 border-b border-gray-100">
                                        <div className="flex items-center gap-4">
                                            <div className="w-14 h-14 rounded-2xl bg-[#EEF0FF] overflow-hidden flex items-center justify-center">
                                                <CategoryIcon className="w-7 h-7 text-blue-400" />
                                            </div>
                                            <div>
                                                <h1 className="font-semibold text-lg text-gray-800">{item.name}</h1>
                                            </div>
                                        </div>
                                    </td>
                                    <td className="px-8 py-6 border-b border-gray-100">
                                        <div>
                                            <h1 className="font-semibold text-gray-800">{item.slug}</h1>
                                        </div>
                                    </td>
                                    <td className="w-[180px] px-8 py-6 border-b border-gray-100">
                                        <span className="px-4 py-2 rounded-full text-sm bg-blue-100 font-semibold">{item.total_used}</span>
                                    </td>
                                </tr>
                            )
                        })}
                    </tbody>
                </table>
            </div>
        </div>
    )
}