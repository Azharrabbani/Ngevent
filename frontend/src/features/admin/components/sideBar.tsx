import React, { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { useLogout } from "../../auth/hooks/useLogout";
import { CgLogOut } from "react-icons/cg";
import { RiCalendarEventFill } from "react-icons/ri";
import { FaRegUser } from "react-icons/fa";
import { IoHomeOutline } from "react-icons/io5";
import { HiMenu } from "react-icons/hi";
import { GoArrowLeft } from "react-icons/go";
import { IoIosArrowDown } from "react-icons/io";
import { HiSwatch } from "react-icons/hi2";
import { useAuth } from "../../../lib/auth";

interface Props {
    children: React.ReactElement;
};

export default function AdminSidebar({children}: Props) {
    const [isMobileOpen, setIsMobileOpen] = useState(false);
    const [isCollapsed, setIsCollapsed] = useState(false);
    const [activeMenu, setActiveMenu] = useState<number | null>(null);
    const location = useLocation();

    const navigate = useNavigate();

    const {logout, loading, error} = useLogout();

    const { user, loading: userLoading } = useAuth()

    const handleLogout = async () => {
        try {
            await logout();
            navigate("/login");
        } catch (err) {
            if (error) {
                return;    
            }
        }
    }

    const menus = [
        { 
            title: "Home", 
            icon: <IoHomeOutline />, 
            path: "/admin/dashboard" 
        },
        {
            title: "Users",
            icon: <FaRegUser />,
            subMenu: true,
            subMenuItems: [
                { title: "Attendee", path: "/admin/attendee-list" },
                { title: "Organizer", path: "/admin/organizer-list"}
            ],
        },
        {
            title: "Events",
            icon: <RiCalendarEventFill />,
            subMenu: true,
            subMenuItems: [
                { title: "Active", path: "/admin/events-active" },
                { title: "Pending", path: "/admin/events-pending"}
            ]
        },
    ];
    
    useEffect(() => {
        menus.forEach((menu, index) => {
            if (menu.subMenu) {
                const isActive = menu.subMenuItems.some(sub => location.pathname === sub.path);

                if (isActive) {
                    setActiveMenu(index);
                }
            }
        });
    }, [location.pathname]);
    
    const logoutMenu = { title: "Logout", icon: <CgLogOut /> };

    return (
        <div className="flex">
            <div className="md:hidden p-4">
                <HiMenu
                    className="text-2xl cursor-pointer"
                    onClick={() => setIsMobileOpen(true)}
                />
            </div>
    
            {isMobileOpen && (
                <div
                    className="fixed inset-0 bg-black/40 z-40 md:hidden"
                    onClick={() => setIsMobileOpen(false)}
                />
            )}
    
            <div className={`
                bg-blue-600 h-screen p-5 pt-8 fixed md:relative z-50 flex flex-col
                ${isCollapsed ? "md:w-20" : "md:w-72"}
                w-72
                ${isMobileOpen ? "left-0" : "-left-full md:left-0"}
                duration-300
            `}>
                <GoArrowLeft
                    className={`bg-white p-1 text-3xl rounded-full absolute -right-3 top-9 border cursor-pointer hidden md:block ${isCollapsed && "rotate-180"}`}
                    onClick={() => setIsCollapsed(!isCollapsed)}
                />

                <GoArrowLeft
                    className="bg-white p-1 text-3xl rounded-full absolute -right-3 top-9 border cursor-pointer md:hidden"
                    onClick={() => setIsMobileOpen(false)}
                />

                <div className="inline-flex">
                    <HiSwatch className="bg-white text-4xl text-blue-500 rounded p-1" />
                    <h1 className={`text-white text-2xl ml-2 duration-300 ${isCollapsed && "hidden"}`}>
                        Ngevent
                    </h1>
                </div>
                
                <div className={`flex flex-col items-center justify-center pt-6 gap-2 duration-300 ${isCollapsed && "hidden"}`}>
                    <div className="bg-white w-20 h-20 rounded-full flex items-center justify-center">
                        <img 
                        className="w-full h-full object-cover rounded-full"
                        src="https://static.vecteezy.com/system/resources/thumbnails/012/210/707/small/worker-employee-businessman-avatar-profile-icon-vector.jpg" 
                        alt="" />
                    </div>
                    {!userLoading && user && (
                        <h1 className="text-center text-white">{user.email}</h1>
                    )}
                </div>

                <div className="flex flex-col justify-between flex-1 overflow-y-auto">
                    <ul className="pt-6">
                        {menus.map((menu, index) => (
                            <div key={index}>
                                <li
                                    className={`group relative text-white flex items-center gap-x-4 p-2 cursor-pointer hover:bg-white/20 rounded-md mt-2 ${location.pathname === menu.path ? "bg-white/30" : "hover:bg-white/20"}`}
                                    onClick={() => {
                                        if (menu.path) {
                                            navigate(menu.path);
                                        }

                                        if (menu.subMenu) {
                                            setActiveMenu(activeMenu === index ? null : index);
                                        }
                                    }}
                                >
                                    <span className="text-2xl">{menu.icon}</span>

                                    <span className={`flex-1 ${isCollapsed && "hidden"}`}>
                                        {menu.title}
                                    </span>
                                    {isCollapsed && (
                                        <span className="absolute left-16 bg-blue-300 text-white text-sm p-2 rounded opacity-0 group-hover:opacity-100 whitespace-nowrap">
                                            {menu.title}
                                        </span>
                                    )}
                                    {menu.subMenu && !isCollapsed && (
                                        <IoIosArrowDown
                                            className={`duration-300 ${activeMenu === index && "rotate-180"}`}
                                        />
                                    )}
                                </li>
                                {menu.subMenu && activeMenu === index && !isCollapsed && (
                                    <ul className="p-2">
                                        {menu.subMenuItems.map((sub, i) => (
                                            <li
                                                key={i}
                                                onClick={() => navigate(sub.path)}
                                                className={`text-white text-sm p-2 pl-10 cursor-pointer rounded-md mt-1
                                                    ${location.pathname === sub.path ? "bg-white/30" : "hover:bg-white/20"}
                                                `}
                                            >
                                                {sub.title}
                                            </li>
                                        ))}
                                    </ul>
                                )}
                            </div>
                        ))}
                    </ul>

                    <ul className="mb-2">
                        <li 
                        className="group relative text-white flex items-center gap-x-4 p-2 cursor-pointer hover:bg-white/50 rounded-md"
                        onClick={() => handleLogout()}
                        >
                            <span className="text-2xl">{logoutMenu.icon}</span>

                            <span className={`${isCollapsed && "hidden"}`}>
                                {loading ? "Logging out..." : logoutMenu.title}
                            </span>
                            {isCollapsed && (
                                <span className="absolute left-16 bg-black text-white text-xs px-2 py-1 rounded opacity-0 group-hover:opacity-100 whitespace-nowrap">
                                    {loading ? "Logging out..." : logoutMenu.title}
                                </span>
                            )}
                        </li>
                    </ul>

                </div>
            </div>

            <div className={`
                flex-1 p-5 mt-8 transition-all duration-300 overflow-y-auto
                ${isCollapsed ? "md:ml-20" : "md:ml-20"}
              `}>
                {children}
            </div>
        </div>
    )
}