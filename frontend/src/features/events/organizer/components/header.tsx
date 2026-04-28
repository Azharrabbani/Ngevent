import { useEffect, useState } from "react";
import { IoIosArrowDown } from "react-icons/io";
import { IoAddOutline } from "react-icons/io5";

export default function Header() {
    const [activeMenu, setActiveMenu] = useState<string | null>(null);
    
    const menus = [
        { title: "Location" },
        {
            title: "Category",
            subMenu: ["Technology", "Business", "Design"]
        },
        {
            title: "Status",
            subMenu: ["draft", "active", "pending", "reject", "done", "cancel"]
        }
    ];

    const handleMenuClick = (title: string) => {
        setActiveMenu(prev => prev === title ? null : title);
    };

    useEffect(() => {
      const handleClickOutside = () => setActiveMenu(null);
      window.addEventListener("click", handleClickOutside);
      return () => window.removeEventListener("click", handleClickOutside);
    }, []);

    return (
        <div className="flex flex-col lg:flex-row justify-between gap-5
                        bg-[#FDFEFF] p-5 md:p-10  w-full"
        >
            <div className="flex flex-col lg:flex-row md:items-center items-center gap-5 md:gap-10">
                <h1 className="text-xl md:text-2xl text-[#0040A1] font-extrabold">Organizer name</h1>
                <div className="flex flex-wrap gap-5 md:gap-10">
                    {menus.map((menu, index) => (
                        <div 
                            key={index} 
                            className="relative"
                        >
                            <h2 
                                onClick={(e) => {
                                    e.stopPropagation();
                                    handleMenuClick(menu.title);
                                }}
                                className="text-sm md:text-lg flex items-center gap-2 cursor-pointer pb-1
                                            text-[#0040A1] font-semibold
                                            hover:border-b-2 hover:border-blue-600"
                            >
                                {menu.title}
                                <IoIosArrowDown className="text-sm"/>
                            </h2>

                            {activeMenu === menu.title && (
                                <div className="absolute left-0 mt-2 w-48 bg-white border border-gray-200 rounded-lg shadow-md z-10 p-3">
                                    {!menu.subMenu && (
                                        <input 
                                            type="text" 
                                            placeholder="Enter location..."
                                            className="w-full border px-2 py-1 rounded"
                                        />
                                    )}

                                    {menu.subMenu && (
                                        <ul className="flex flex-col gap-2">
                                            {menu.subMenu.map((item, i) => (
                                                <li
                                                    key={i}
                                                    className="cursor-pointer hover:bg-gray-100 px-2 py-1 rounded"
                                                >
                                                    {item}
                                                </li>
                                            ))}
                                        </ul>
                                    )}

                                </div>
                            )}

                        </div>
                    ))}                    
                </div>
            </div>

            <button className="flex items-center justify-center gap-2
                                bg-[#0040A1] rounded-full hover:shadow-xl 
                                transition-all duration-200
                                px-6 py-2 md:px-8 md:py-3
                                text-white text-sm md:text-lg"
            >                    
                <IoAddOutline/>
                Create Event
            </button>
        </div>
    )
}