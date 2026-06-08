import { useMemo, useRef } from "react";
import {
    FaChevronLeft,
    FaChevronRight,
} from "react-icons/fa";
import useCategorySlider from "../../utils/useCategorySlider";
import { useListCategories } from "../../../../categories/hooks/useListCategories";
import CategoryItem from "./categoryItem";

interface Props {
    selectedCategory?: number[];
    onChange: (
        categoryIds?: number[]
    ) => void;
}

export default function CategorySlider({
    selectedCategory,
    onChange,
}: Props) {
    const sliderRef =
        useRef<HTMLDivElement>(null);

    const {
        canScrollLeft,
        canScrollRight,
        scrollLeft,
        scrollRight,
    } = useCategorySlider(sliderRef);

    const { data: categories = [] } =
        useListCategories();

    const categoryItems = useMemo(
        () => [
            {
                id: 0,
                name: "All Events",
            },
            ...categories,
        ],
        [categories]
    );

    return (
        <div className="relative mt-6">
            {/* Left scroll button */}
            {canScrollLeft && (
                <button
                    onClick={scrollLeft}
                    className="
                        absolute left-0
                        top-1/2 -translate-y-1/2
                        z-10
                        w-8 h-8
                        bg-white border border-slate-200
                        rounded-full shadow-md
                        flex items-center justify-center
                        text-slate-500 hover:text-blue-600
                        transition
                    "
                >
                    <FaChevronLeft size={12} />
                </button>
            )}

            {/* Right scroll button */}
            {canScrollRight && (
                <button
                    onClick={scrollRight}
                    className="
                        absolute right-0
                        top-1/2 -translate-y-1/2
                        z-10
                        w-8 h-8
                        bg-white border border-slate-200
                        rounded-full shadow-md
                        flex items-center justify-center
                        text-slate-500 hover:text-blue-600
                        transition
                    "
                >
                    <FaChevronRight size={12} />
                </button>
            )}

            {/* Scrollable container */}
            <div
                ref={sliderRef}
                className="
                    flex gap-2
                    overflow-x-auto
                    scrollbar-hide
                    px-10
                    py-1
                "
            >
                {categoryItems.map(
                    (category) => (
                        <CategoryItem
                            key={category.id}
                            label={category.name}
                            active={
                                category.id === 0
                                    ? !selectedCategory
                                        ?.length
                                    : selectedCategory?.includes(
                                        category.id
                                    )
                            }
                            onClick={() => {
                                if (
                                    category.id ===
                                    0
                                ) {
                                    onChange(
                                        undefined
                                    );

                                    return;
                                }

                                onChange([
                                    category.id,
                                ]);
                            }}
                        />
                    )
                )}
            </div>
        </div>
    );
}
