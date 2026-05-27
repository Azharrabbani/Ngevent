import { useState } from "react";
import Button from "../../../../components/Button";
import { PlusIcon, SearchIcon } from "../../../../components/icon";
import Input from "../../../../components/input";
import CreateCategorySlideOver from "./createCategorySlideOver";

interface Props {
    search?: string;
    setSearch?: (val: string | undefined) => void;
}

export default function CategoriesHeader({ search, setSearch }: Props) {
    const [isSlideOverOpen, setIsSlideOverOpen] = useState(false);

    return (
        <>
            <div className="flex flex-col gap-y-5 p-6 md:flex-row md:items-center md:justify-between">
                <h1 className="text-3xl font-bold text-center md:text-left">Categories</h1>

                <div className="flex flex-col md:flex-row gap-4 items-center">
                    <div className="relative w-full md:w-[280px]">
                        <Input
                            leftIcon={<SearchIcon />}
                            type="text"
                            placeholder="Search categories"
                            value={search}
                            onChange={(e) => setSearch?.(e.target.value)}
                            className="w-full bg-white pl-10 pr-4 py-3 rounded-lg border border-gray-300 outline-none focus:ring-2 focus:ring-blue-500 text-sm"
                        />
                    </div>

                    <Button
                        className="px-6 py-3 rounded-xl bg-[#0066FF]"
                        onClick={() => setIsSlideOverOpen(true)}
                    >
                        <div className="flex justify-center items-center gap-2">
                            <PlusIcon />
                            <p>Create New Category</p>
                        </div>
                    </Button>
                </div>
            </div>

            <CreateCategorySlideOver
                isOpen={isSlideOverOpen}
                onClose={() => setIsSlideOverOpen(false)}
            />
        </>
    );
}