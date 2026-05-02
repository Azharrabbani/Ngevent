import { IoCheckmark, IoChevronDown, IoClose } from "react-icons/io5";

interface Props {
    selectedItems: any[];
    toogleItem: (item: number) => void;
    isItemOpen: boolean;
    itemsLoading: boolean;
    setIsItemOpen: (item: any) => void;
    items: any[] | undefined;
    placeholder: string;
};

export default function SelectItems({
    selectedItems,
    toogleItem,
    isItemOpen,
    itemsLoading,
    setIsItemOpen,
    items,
    placeholder
}: Props) {
    return (
        <div className="relative min-w-0 w-full">
            <label className="text-xs sm:text-sm text-[#003B95] font-medium mb-1 block">
                Category
            </label>

            {/* Trigger for dropdown */}
            <div
                className="min-h-12 w-full bg-gray-200 rounded-xl px-3 py-2
                            flex items-start gap-2 cursor-pointer"
                onClick={() => setIsItemOpen(!isItemOpen)}
            >
                {selectedItems.length > 0 ? (
                    <div className="flex-1 min-w-0 flex flex-wrap gap-2 max-h-20 overflow-y-auto">
                         {selectedItems.map((id) => {
                            const item = items?.find((item) => item.id === id);
                            return (
                                <span
                                    key={id}
                                    className="max-w-full px-3 py-1 bg-[#003B95] text-white rounded-lg text-sm flex items-center gap-2"
                                >
                                    <span className="truncate">
                                        {item.name}
                                    </span>
                                    <button
                                        type="button"
                                        onClick={(e) => {
                                            e.stopPropagation();
                                            toogleItem(id)
                                        }}
                                    >
                                        <IoClose size={16}/>
                                    </button>
                                </span>
                            )
                        })}
                    </div>
                ) : (
                    <span className="flex-1 text-gray-500 text-sm">
                        {placeholder}
                    </span>
                )}
                <IoChevronDown className={`ml-auto shrink-0 mt-1 transition ${
                  isItemOpen ? "rotate-180" : ""
                }`}/>
            </div>
            
            {isItemOpen && (
                <div className="absolute z-50 mt-2 w-full bg-white border border-gray-200 rounded-xl shadow-lg max-h-72 overflow-y-auto">
                    {itemsLoading ? (
                        <span className="w-full px-4 py-4 flex items-center justify-between text-left">Loading...</span>
                    ) : !items?.length ? (
                        <span className="w-full px-4 py-4 flex items-center justify-between text-left">Item not found</span>
                    ) : (
                        items?.map((item) => (
                            <button
                                key={item.id}
                                type="button"
                                onClick={() => toogleItem(item.id)}
                                className="w-full px-4 py-4 flex items-center justify-between hover:bg-gray-50 text-left"
                            >   
                                <span>{item.name}</span>
                                {selectedItems.includes(item.id) && (
                                    <IoCheckmark
                                        className="text-[#003B95]"
                                        size={20}
                                    />
                                )}
                            </button>
                        ))
                    )
                    }
                </div>
            )}
        </div>
    )
}