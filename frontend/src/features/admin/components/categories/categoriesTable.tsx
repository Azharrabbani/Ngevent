import { useEffect, useRef, useState } from "react";
import { CategoryIcon } from "../../../../components/icon";
import type { PaginatedData } from "../../../../types/apiResponse";
import type { categoriesPaginatedResp } from "../../../categories/types/categoryResponse";
import { useUpdateCategory } from "../../../categories/hooks/useUpdateCategory";
import Input from "../../../../components/input";

interface Props {
    data: PaginatedData<categoriesPaginatedResp>;
}

interface EditState {
    id: number | string;
    value: string;
}

export default function CategoriesTable({ data }: Props) {
    const [editing, setEditing] = useState<EditState | null>(null);
    const inputRef = useRef<HTMLInputElement>(null);
    const wrapperRef = useRef<HTMLDivElement>(null);

    const { mutate: updateCategory, isPending } = useUpdateCategory();

    useEffect(() => {
        if (editing?.id) {
            inputRef.current?.focus();
            inputRef.current?.select();
        }
    }, [editing?.id]);

    const handleDoubleClick = (item: categoriesPaginatedResp) => {
        setEditing({ id: item.id, value: item.name });
    };

    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (
                editing &&
                wrapperRef.current &&
                !wrapperRef.current.contains(event.target as Node)
            ) {
                setEditing(null);
            }
        };

        document.addEventListener("mousedown", handleClickOutside);

        return () => {
            document.removeEventListener(
                "mousedown",
                handleClickOutside
            );
        };
    }, [editing]);

    const handleSave = () => {
        if (!editing) return;
        const trimmed = editing.value.trim();

        const original = data.rows.find((r) => r.id === editing.id);
        if (!trimmed || trimmed === original?.name) {
            setEditing(null);
            return;
        }

        updateCategory(
            { id: editing.id, name: trimmed },
            { onSettled: () => setEditing(null) }
        );
    };

    const handleKeyDown = (
        e: React.KeyboardEvent<HTMLInputElement>
    ) => {
        e.stopPropagation();

        if (e.key === "Enter") {
            e.preventDefault();
            handleSave();
        }

        if (e.key === "Escape") {
            e.preventDefault();
            setEditing(null);
        }
    };
    return (
        <div className="hidden sm:flex justify-center w-full">
            <div className="w-full lg:max-w-8xl overflow-hidden border-2 border-[#C2C6D8] rounded-xl">
                <table className="w-full table-auto">
                    <thead className="bg-[#EFF4FF]">
                        <tr>
                            <th className="px-8 py-5 text-left text-md font-semibold text-gray-600 border-b-2 border-[#C2C6D8]">
                                CATEGORY NAME
                                <span className="ml-2 text-xs font-normal text-gray-400">(double-click to edit)</span>
                            </th>
                            <th className="px-8 py-5 text-left text-md font-semibold text-gray-600 border-b-2 border-[#C2C6D8]">SLUG</th>
                            <th className="w-[180px] px-8 py-5 text-left text-md font-semibold text-gray-600 border-b-2 border-[#C2C6D8]">
                                TOTAL USED
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        {data.rows.map((item) => {
                            const isEditing = editing?.id === item.id;
                            const isSaving = isPending && isEditing;

                            return (
                                <tr
                                    key={item.id}
                                    className="hover:bg-gray-100 transition-colors duration-200 border-b-2 border-[#C2C6D8] last:border-0"
                                    onDoubleClick={() => !editing && handleDoubleClick(item)}
                                >
                                    <td className="px-8 py-6 border-b border-gray-100">
                                        <div className="flex items-center gap-4">
                                            <div className="w-14 h-14 rounded-2xl bg-[#EEF0FF] overflow-hidden flex items-center justify-center shrink-0">
                                                <CategoryIcon className="w-7 h-7 text-blue-400" />
                                            </div>

                                            {isEditing ? (
                                                <div
                                                    ref={wrapperRef}
                                                    className="flex items-center gap-2 flex-1"
                                                >
                                                    <Input
                                                        ref={inputRef}
                                                        type="text"
                                                        value={editing.value}
                                                        onChange={(e) =>
                                                            setEditing((prev) =>
                                                                prev ? { ...prev, value: e.target.value } : null
                                                            )
                                                        }
                                                        onKeyDown={handleKeyDown}
                                                        disabled={isSaving}
                                                        className="flex-1 px-3 py-1.5 rounded-lg border-2 bg-white border-blue-500 text-base font-semibold text-gray-800 outline-none focus:ring-2 focus:ring-blue-200 bg-white disabled:opacity-60"
                                                    />
                                                    <span className="text-xs text-gray-400 whitespace-nowrap">
                                                        Enter ↵ or Esc
                                                    </span>
                                                </div>
                                            ) : (
                                                <div className="cursor-text select-none">
                                                    <h1 className="font-semibold text-lg text-gray-800">
                                                        {isSaving ? (
                                                            <span className="opacity-50">{item.name}</span>
                                                        ) : (
                                                            item.name
                                                        )}
                                                    </h1>
                                                </div>
                                            )}
                                        </div>
                                    </td>

                                    <td className="px-8 py-6 border-b border-gray-100">
                                        <h1 className="font-semibold text-gray-800">{item.slug}</h1>
                                    </td>

                                    <td className="w-[180px] px-8 py-6 border-b border-gray-100">
                                        <span className="px-4 py-2 rounded-full text-sm bg-blue-100 font-semibold">
                                            {item.total_used}
                                        </span>
                                    </td>
                                </tr>
                            );
                        })}
                    </tbody>
                </table>
            </div>
        </div>
    );
}