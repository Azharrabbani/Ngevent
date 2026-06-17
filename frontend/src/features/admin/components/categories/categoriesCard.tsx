import { useEffect, useRef, useState } from "react";
import { CategoryIcon } from "../../../../components/icon";
import type { PaginatedData } from "../../../../types/apiResponse";
import type { categoriesPaginatedResp } from "../../../categories/types/categoryResponse";
import { useUpdateCategory } from "../../../categories/hooks/useUpdateCategory";
import Input from "../../../../components/input";
import { mapValidationErrors } from "../../../../utils/validation";

interface Props {
    data: PaginatedData<categoriesPaginatedResp>;
}

interface EditState {
    id: number | string;
    value: string;
}

export default function CategoriesCard({ data }: Props) {
    const [editing, setEditing] = useState<EditState | null>(null);
    const inputRef = useRef<HTMLInputElement>(null);
    const wrapperRef = useRef<HTMLDivElement>(null);

    const { mutate: updateCategory, isPending, error } = useUpdateCategory();

    const validationErrors = mapValidationErrors(error);

    useEffect(() => {
        if (editing?.id) {
            inputRef.current?.focus();
            inputRef.current?.select();
        }
    }, [editing?.id]);

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

    const handleDoubleClick = (item: categoriesPaginatedResp) => {
        setEditing({ id: item.id, value: item.name });
    };

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
        <div className="sm:hidden p-4 space-y-4 max-h-150 sm:max-h-125 overflow-y-auto">
            {data.rows.map((item) => {
                const isEditing = editing?.id === item.id;
                const isSaving = isPending && isEditing;

                return (
                    <div
                        key={item.id}
                        className="hover:bg-gray-50 transition-colors duration-200 border border-gray-100 rounded-2xl p-4 shadow-sm bg-white"
                        onDoubleClick={() => !editing && handleDoubleClick(item)}
                    >
                        <div className="flex items-start gap-4">
                            <div className="w-14 h-14 rounded-2xl bg-[#EEF0FF] flex items-center justify-center shrink-0">
                                <CategoryIcon className="w-7 h-7 text-blue-400" />
                            </div>

                            <div className="flex-1 min-w-0">
                                <div className="flex items-start justify-between gap-3">
                                    {isEditing ? (
                                        <div
                                            ref={wrapperRef}
                                            className="flex-1 flex flex-col gap-1"
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
                                                className="w-full px-3 py-1.5 bg-white rounded-lg border-2 border-blue-500 text-base font-semibold text-gray-800 outline-none focus:ring-2 focus:ring-blue-200 bg-white disabled:opacity-60"
                                            />
                                            {validationErrors?.message?.name && (
                                                <p className="text-red-500 text-xs">
                                                    {validationErrors.message.name}
                                                </p>
                                            )}
                                            <span className="text-xs text-gray-400">
                                                Tap Enter to save, Esc to cancel
                                            </span>
                                        </div>
                                    ) : (
                                        <h1 className="font-semibold text-base text-gray-800 leading-snug select-none cursor-text">
                                            {item.name}
                                        </h1>
                                    )}

                                    <span className="px-4 py-2 rounded-full text-sm bg-blue-100 font-semibold whitespace-nowrap shrink-0">
                                        {item.total_used}
                                    </span>
                                </div>

                                <div className="mt-4">
                                    <p className="text-xs text-gray-400">Slug</p>
                                    <h2 className="text-sm font-medium text-gray-700">{item.slug}</h2>
                                </div>

                                {!isEditing && (
                                    <p className="mt-2 text-xs text-gray-300">Double-tap to edit name</p>
                                )}
                            </div>
                        </div>
                    </div>
                );
            })}
        </div>
    );
}