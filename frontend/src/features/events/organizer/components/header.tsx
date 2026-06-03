import { useEffect, useState, type Dispatch, type SetStateAction } from "react";
import { IoIosArrowDown, IoIosSearch } from "react-icons/io";
import { IoAddOutline } from "react-icons/io5";
import type { categoriesResp } from "../../../categories/types/categoryResponse";
import { useNavigate } from "react-router-dom";
import type { DateFilterType } from "../../../../utils/dateFilter";
import DateFilterDropdown from "../../../../components/dateFilter/dateFilterDropdown";

interface Props {
    organizerName: string | undefined;
    location: string | undefined;
    setLocation: (val: string | undefined) => void;
    event: string | undefined;
    categories: categoriesResp[] | undefined;
    categoriesLoading: boolean;
    selectedCategories: number[];
    status?: string | undefined;
    setStatus?: Dispatch<SetStateAction<string | undefined>>;
    setSelectedCategories: Dispatch<SetStateAction<number[]>>;
    dateFilterType: DateFilterType;
    setDateFilterType:
    Dispatch<SetStateAction<DateFilterType>>;

    selectedDate: Date | null;
    setSelectedDate:
    Dispatch<SetStateAction<Date | null>>;

    selectedMonth?: number;
    setSelectedMonth:
    Dispatch<SetStateAction<number | undefined>>;

    selectedYear?: number;
    setSelectedYear:
    Dispatch<SetStateAction<number | undefined>>;
    setEvent: (val: string | undefined) => void;
    onSearch: () => void;
    toggleStatus: boolean
};

export default function Header({
    organizerName,
    location,
    setLocation,
    event,
    categories,
    categoriesLoading,
    selectedCategories,
    status,
    setStatus,
    dateFilterType,
    setDateFilterType,
    selectedDate,
    setSelectedDate,
    selectedMonth,
    setSelectedMonth,
    selectedYear,
    setSelectedYear,
    setSelectedCategories,
    setEvent,
    toggleStatus,
}: Props) {
    const navigate = useNavigate();

    const [activeMenu, setActiveMenu] = useState<string | null>(null);

    const menus = [
        { title: "Location" },
        { title: "Category" },
        ...(toggleStatus
            ? [
                {
                    title: "Status",
                    subMenu: ["active", "pending", "rejected", "done"],
                }
            ]
            : []),
        { title: "Date" },
    ];

    const handleMenuClick = (title: string) => {
        setActiveMenu(prev => prev === title ? null : title);
    };

    const toggleCategory = (id: number) => {
        setSelectedCategories(prev =>
            prev.includes(id)
                ? prev.filter(c => c !== id)
                : [...prev, id]
        );
    };

    useEffect(() => {
        const handleClickOutside = () => setActiveMenu(null);
        window.addEventListener("click", handleClickOutside);
        return () => window.removeEventListener("click", handleClickOutside);
    }, []);

    return (
        <div className="flex flex-col xl:flex-row justify-between gap-5 bg-[#FDFEFF] p-5 md:p-10 w-full">

            {/* Overlay (Mobile Only) */}
            {activeMenu && (
                <div
                    className="fixed inset-0 bg-black/20 z-40 md:hidden"
                    onClick={() => setActiveMenu(null)}
                />
            )}

            <div className="flex flex-col xl:flex-row md:items-center items-center gap-5 md:gap-10">
                <h1 className="text-xl md:text-2xl text-[#0040A1] font-extrabold">
                    {organizerName ? organizerName : "Ngevent"}
                </h1>

                <div className="flex flex-wrap gap-5 md:gap-10">
                    {menus.map((menu, index) => (
                        <div key={index} className="relative">

                            {/* MENU BUTTON */}
                            <h2
                                onClick={(e) => {
                                    e.stopPropagation();
                                    handleMenuClick(menu.title);
                                }}
                                className="text-sm md:text-lg flex items-center gap-2 cursor-pointer pb-1 text-[#0040A1] font-semibold hover:border-b-2 hover:border-blue-600"
                            >
                                {menu.title}
                                <IoIosArrowDown />
                            </h2>

                            {/* DROPDOWN */}
                            {activeMenu === menu.title && (
                                <div
                                    onClick={(e) => e.stopPropagation()}
                                    className="
                                        fixed bottom-0 left-0 right-0 w-full z-50
                                        md:absolute md:bottom-auto md:w-56

                                        bg-white border border-gray-200 
                                        rounded-t-2xl md:rounded-lg 
                                        shadow-lg p-4
                                        max-h-[60vh] overflow-y-auto
                                    "
                                >

                                    {/* MOBILE HEADER */}
                                    <div className="md:hidden mb-3 text-center font-semibold">
                                        {menu.title}
                                    </div>

                                    {/* LOCATION */}
                                    {menu.title === "Location" && (
                                        <form
                                            onSubmit={(e) => {
                                                e.preventDefault();
                                            }}
                                            className="relative"
                                        >
                                            <input
                                                type="text"
                                                placeholder="Enter Location..."
                                                value={location}
                                                onChange={(e) => setLocation(e.target.value)}
                                                className="w-full border px-3 py-2 rounded-lg pr-10 focus:ring-2 focus:ring-blue-400"
                                            />
                                            <IoIosSearch className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400" />
                                        </form>
                                    )}

                                    {/* DATE */}
                                    {menu.title === "Date" && (
                                        <DateFilterDropdown
                                            dateFilterType={dateFilterType}
                                            setDateFilterType={setDateFilterType}
                                            selectedDate={selectedDate}
                                            setSelectedDate={setSelectedDate}
                                            selectedMonth={selectedMonth}
                                            setSelectedMonth={setSelectedMonth}
                                            selectedYear={selectedYear}
                                            setSelectedYear={setSelectedYear}
                                        />
                                    )}

                                    {/* CATEGORY */}
                                    {menu.title === "Category" && (
                                        <div className="flex flex-col gap-2">
                                            {categoriesLoading ? (
                                                <p className="text-sm text-gray-400 text-center">
                                                    Loading categories...
                                                </p>
                                            ) : !categories?.length ? (
                                                <p className="text-sm text-gray-400 text-center">
                                                    No categories
                                                </p>
                                            ) : (
                                                categories.map(cat => {
                                                    const isSelected = selectedCategories.includes(cat.id);

                                                    return (
                                                        <div
                                                            key={cat.id}
                                                            onClick={() => toggleCategory(cat.id)}
                                                            className={`flex justify-between items-center px-2 py-1 rounded cursor-pointer
                                                                ${isSelected
                                                                    ? "bg-blue-100 text-blue-600 font-medium"
                                                                    : "hover:bg-gray-100"
                                                                }`}
                                                        >
                                                            <span>{cat.name}</span>
                                                            {isSelected && "✔"}
                                                        </div>
                                                    );
                                                })
                                            )}
                                        </div>
                                    )}

                                    {menu.title === "Status" && toggleStatus && (
                                        <div className="flex flex-col gap-2">
                                            {menu.subMenu?.map((item, i) => {
                                                const isSelected = status === item;

                                                return (
                                                    <div
                                                        key={i}
                                                        onClick={(e) => {
                                                            e.stopPropagation();
                                                            setStatus?.(prev => prev === item ? undefined : item);
                                                        }}
                                                        className={`flex justify-between items-center px-2 py-1 rounded cursor-pointer
                                                            ${isSelected
                                                                ? "bg-blue-100 text-blue-600 font-medium"
                                                                : "hover:bg-gray-100"
                                                            }`}
                                                    >
                                                        <span>{item}</span>
                                                        {isSelected && "✔"}
                                                    </div>
                                                );
                                            })}
                                        </div>
                                    )}

                                </div>
                            )}
                        </div>
                    ))}
                </div>
            </div>

            <div className="flex flex-col xl:flex-row items-center gap-3 w-full lg:w-auto">
                <form
                    className="relative w-full xl:w-72"
                    onSubmit={(e) => {
                        e.preventDefault();
                    }}
                >
                    <IoIosSearch className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
                    <input
                        type="text"
                        placeholder="Search events..."
                        value={event}
                        onChange={(e) => setEvent(e.target.value)}
                        className="w-full border px-3 py-2 pl-10 rounded-full focus:ring-2 focus:ring-blue-400"
                    />
                </form>

                <button
                    className="flex items-center justify-center gap-2 
                                bg-[#0040A1] rounded-full hover:shadow-xl transition-all duration-200 
                                px-6 py-2 md:px-8 md:py-3 text-white text-sm md:text-lg"
                    onClick={() => navigate("/organizer/event/new")}
                >
                    <IoAddOutline />
                    Create Event
                </button>
            </div>
        </div>
    );
}