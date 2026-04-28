import { useState } from "react";
import { CgProfile } from "react-icons/cg";
import { FaRegCircleCheck } from "react-icons/fa6";
import { HiMenu } from "react-icons/hi";
import { IoMdSave } from "react-icons/io";
import { MdOutlineCancel } from "react-icons/md";
import { useNavigate } from "react-router-dom";

interface Props {
    children: React.ReactElement;
};

export default function Sidebar({ children }: Props) {
    const navigate = useNavigate();

    const [isOpen, setIsOpen] = useState(false);

    const menus = [
        { 
            title: "Cancel", 
            icon: <MdOutlineCancel />, 
            path: "/organizer/cancel-event" 
        },
        {
            title: "Check",
            icon: <FaRegCircleCheck />,
            subMenu: true,
            path: "/organizer/active-event"
        },
        {
            title: "Save",
            icon: <IoMdSave />,
            subMenu: true,
            path: "/organizer/draft-event"
        },
        {
            title: "Profile",
            icon: <CgProfile />,
            subMenu: true,
            path: "/profile"
        },
    ];

    return (
        <div className="flex">

            <button
                className="md:hidden fixed top-4 left-4 z-50 bg-white p-2 rounded-lg shadow"
                onClick={() => setIsOpen(true)}
            >
                <HiMenu className="text-2xl"/>
            </button>

            {isOpen && (
                <div 
                    className="fixed inset-0 bg-black/40 z-40 md:hidden"
                    onClick={() => setIsOpen(false)}
                />

            )}

            <div 
                className={`
                    fixed md:fixed top-0 left-0 h-screen z-50
                    w-20 bg-[#EDEDF8] p-5 pt-12 flex flex-col gap-10
                    transform transition-transform duration-300
                    ${isOpen ? "translate-x-0" : "-translate-x-full"}
                    md:translate-x-0
                    `}
            >
                <img 
                    src="https://static.vecteezy.com/system/resources/thumbnails/040/890/255/small/ai-generated-empty-wooden-table-on-the-natural-background-for-product-display-free-photo.jpg"
                    alt="" 
                    className="w-15 h-10 object-cover rounded-full mx-auto"
                />

                <div className="flex flex-col gap-4">
                    {menus.map((menu, index) => (
                        <div key={index}>
                            <li 
                                className="flex flex-col px-3 py-4 rounded-2xl
                                            items-center gap-1 text-gray-600
                                            hover:bg-[#0056D2] hover:text-white hover:shadow-xl hover:shadow-blue-500 transition-all duration-200
                                            cursor-pointer"
                                onClick={() => {
                                    navigate(menu.path)
                                    setIsOpen(false)
                                }}
                            >
                                <span className="text-xl">{menu.icon}</span>
                                <h2 className="text-xs">{menu.title}</h2>
                            </li>
                        </div>
                    ))}
                </div>
            </div>

            <div className="w-full md:ml-20 min-h-screen overflow-y-auto">
                {children}
            </div>
        </div>
    )
}