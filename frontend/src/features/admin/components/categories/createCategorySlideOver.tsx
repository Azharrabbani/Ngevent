import { useEffect, useRef, useState } from "react";
import Button from "../../../../components/Button";
import Input from "../../../../components/input";
import { useCreateCategory } from "../../../categories/hooks/useCreateCategory";
import { CrossIcon } from "../../../../components/icon";
import { CreateSlug } from "../../utils/slug";
import Slider from "../slider";

interface Props {
    isOpen: boolean;
    onClose: () => void;
}

export default function CreateCategorySlideOver({ isOpen, onClose }: Props) {
    const [name, setName] = useState<string>("");
    const inputRef = useRef<HTMLInputElement>(null);

    const slug = CreateSlug(name);

    const { mutate: createCategory, isPending } = useCreateCategory(() => {
        setName("");
        onClose();
    });

    useEffect(() => {
        if (isOpen) {
            setTimeout(() => inputRef.current?.focus(), 150);
        } else {
            setName("");
        }
    }, [isOpen]);

    useEffect(() => {
        const handleKey = (e: KeyboardEvent) => {
            if (e.key === "Escape" && isOpen) onClose();
        };
        document.addEventListener("keydown", handleKey);
        return () => document.removeEventListener("keydown", handleKey);
    }, [isOpen, onClose]);

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (!name.trim()) return;
        createCategory(name.trim());
    };

    return (
        <Slider isOpen={isOpen} onClose={onClose}>
            <>
                <div className="flex items-center justify-between px-6 py-5 border-b border-gray-100">
                    <div>
                        <h2 className="text-xl font-bold text-gray-900">New Category</h2>
                        <p className="text-sm text-gray-500 mt-0.5">
                            Add a new category to the system
                        </p>
                    </div>
                    <button
                        onClick={onClose}
                        className="w-9 h-9 flex items-center justify-center rounded-xl hover:bg-gray-100 transition-colors text-gray-500 hover:text-gray-800"
                    >
                        <CrossIcon className="w-5 h-5" />
                    </button>
                </div>

                <form onSubmit={handleSubmit} className="flex flex-col flex-1 px-6 py-6 gap-6">
                    <div className="flex flex-col gap-1.5">
                        <label className="text-sm font-semibold text-gray-700">
                            Category Name <span className="text-red-500">*</span>
                        </label>
                        <Input
                            ref={inputRef}
                            type="text"
                            placeholder="e.g. Technology"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            className="w-full bg-white px-4 py-3 rounded-lg border border-gray-300 outline-none focus:ring-2 focus:ring-blue-500 text-sm"
                        />
                    </div>

                    <div className="flex flex-col gap-1.5">
                        <label className="text-sm font-semibold text-gray-700">
                            Slug Preview
                        </label>
                        <div className="px-4 py-3 rounded-lg bg-gray-50 border border-dashed border-gray-300">
                            <span className="text-sm text-gray-500 font-mono">
                                {slug || (
                                    <span className="text-gray-400 italic">
                                        will be generated from name...
                                    </span>
                                )}
                            </span>
                        </div>
                        <p className="text-xs text-gray-400">
                            Auto-generated from the category name. Cannot be edited manually.
                        </p>
                    </div>

                    <div className="flex-1" />

                    <div className="flex gap-3 pt-4 border-t border-gray-100">
                        <button
                            type="button"
                            onClick={onClose}
                            disabled={isPending}
                            className="flex-1 px-4 py-3 rounded-xl border border-gray-200 text-gray-700 font-semibold text-sm hover:bg-gray-50 transition-colors disabled:opacity-50"
                        >
                            Cancel
                        </button>
                        <Button
                            type="submit"
                            disabled={!name.trim() || isPending}
                            className="flex-1 px-4 py-3 rounded-xl bg-[#0066FF] text-white font-semibold text-sm disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                        >
                            {isPending ? (
                                <>
                                    <span className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                                    Creating...
                                </>
                            ) : (
                                "Create Category"
                            )}
                        </Button>
                    </div>
                </form>
            </>
        </Slider>
    );
}