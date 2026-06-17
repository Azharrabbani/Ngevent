import React, { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { useLogout } from "../../auth/hooks/useLogout";
import { CgLogOut } from "react-icons/cg";
import { HiMenu } from "react-icons/hi";
import { GoArrowLeft } from "react-icons/go";
import { IoIosArrowDown } from "react-icons/io";
import { useAuth } from "../../../lib/auth";
import { CategoryIcon, CloseIcon, DashboardIcon, EventIcon, HomeIcon, UserIcon } from "../../../components/icon";

interface Props {
    children: React.ReactElement;
}

const LG_BREAKPOINT = 1024;

export default function AdminSidebar({ children }: Props) {
    const [isMobileOpen, setIsMobileOpen] = useState(false);
    const [isCollapsed, setIsCollapsed] = useState(false);
    const [activeMenu, setActiveMenu] = useState<number | null>(null);
    const location = useLocation();
    const navigate = useNavigate();

    const { mutateAsync: logout, isPending, isError } = useLogout();
    const { user, loading: userLoading } = useAuth();

    useEffect(() => {
        const handleResize = () => {
            if (window.innerWidth >= LG_BREAKPOINT) {
                setIsMobileOpen(false);
            }
        };
        window.addEventListener("resize", handleResize);
        return () => window.removeEventListener("resize", handleResize);
    }, []);

    useEffect(() => {
        setIsMobileOpen(false);
    }, [location.pathname]);

    const handleLogout = async () => {
        try {
            await logout();
            navigate("/");
        } catch (err) {
            if (isError) return;
        }
    };

    const menus = [
        {
            title: "Home",
            icon: <HomeIcon />,
            path: "/admin/dashboard",
        },
        {
            title: "Users",
            icon: <UserIcon />,
            subMenu: true,
            subMenuItems: [
                { title: "Admin", path: "/admin/admin-list" },
                { title: "Organizer", path: "/admin/organizer-list" },
            ],
        },
        {
            title: "Events",
            icon: <EventIcon />,
            subMenu: true,
            subMenuItems: [
                { title: "Active", path: "/admin/events/active" },
                { title: "Pending", path: "/admin/events/pending" },
                { title: "Done", path: "/admin/events/done" },
                { title: "Rejected", path: "/admin/events/rejected" },
            ],
        },
        {
            title: "Categories",
            icon: <CategoryIcon />,
            path: "/admin/categories",
        },
        {
            title: "Dashboard",
            icon: <DashboardIcon />,
            path: "/",
        },
    ];

    useEffect(() => {
        menus.forEach((menu, index) => {
            if (menu.subMenu) {
                const isActive = menu.subMenuItems.some(
                    (sub) => location.pathname === sub.path
                );
                if (isActive) setActiveMenu(index);
            }
        });
    }, [location.pathname]);

    const logoutMenu = { title: "Logout", icon: <CgLogOut /> };

    return (
        <div className="flex h-screen overflow-hidden bg-gray-50">

            <div className="lg:hidden fixed top-4 left-4 z-50">
                <button
                    className="p-2 rounded-lg bg-blue-600 text-white shadow-md"
                    onClick={() => setIsMobileOpen(true)}
                >
                    <HiMenu className="text-xl" />
                </button>
            </div>

            {isMobileOpen && (
                <div
                    className="fixed inset-0 bg-black/50 z-40 lg:hidden"
                    onClick={() => setIsMobileOpen(false)}
                />
            )}

            <aside
                className={`
                    bg-blue-600 h-screen flex flex-col shrink-0
                    transition-all duration-300

                    fixed top-0 left-0 z-50 w-72
                    lg:sticky lg:top-0
                    ${isCollapsed ? "lg:w-20" : "lg:w-72"}
                    ${isMobileOpen ? "translate-x-0" : "-translate-x-full lg:translate-x-0"}
                `}
            >
                <GoArrowLeft
                    className={`
                        bg-white p-1 text-3xl rounded-full
                        absolute -right-3 top-9 border cursor-pointer
                        hidden lg:block
                        transition-transform duration-300
                        ${isCollapsed ? "rotate-180" : ""}
                    `}
                    onClick={() => setIsCollapsed(!isCollapsed)}
                />

                {/* Mobile close button */}
                <button
                    className="lg:hidden absolute top-4 right-4 text-white/80 hover:text-white p-1 rounded-lg hover:bg-white/20 transition-colors"
                    onClick={() => setIsMobileOpen(false)}
                    aria-label="Close sidebar"
                >
                    <CloseIcon className="w-5 h-5" />
                </button>

                <div className="inline-flex items-center p-5 pt-8">
                    <div className="bg-white rounded-full shrink-0">
                        <img
                            src="/initial_logo.png"
                            alt="Ngevent"
                            className="w-10 h-10 object-contain"
                        />
                    </div>
                    <h1 className={`text-white text-2xl ml-2 duration-300 ${isCollapsed ? "hidden" : ""}`}>
                        Ngevent
                    </h1>
                </div>

                <div className={`flex flex-col items-center justify-center gap-2 px-5 duration-300 ${isCollapsed ? "hidden" : ""}`}>
                    <div className="bg-white w-20 h-20 rounded-full flex items-center justify-center overflow-hidden">
                        <img
                            className="w-full h-full object-cover rounded-full"
                            src="https://static.vecteezy.com/system/resources/thumbnails/012/210/707/small/worker-employee-businessman-avatar-profile-icon-vector.jpg"
                            alt=""
                        />
                    </div>
                    {!userLoading && user && (
                        <h1 className="text-center text-white text-sm break-all px-2">
                            {user.email}
                        </h1>
                    )}
                </div>

                {/* Nav */}
                <div className="flex flex-col justify-between flex-1 overflow-y-auto px-5 pb-4">
                    <ul className="pt-6 space-y-1">
                        {menus.map((menu, index) => (
                            <div key={index}>
                                <li
                                    className={`
                                        group relative text-white flex items-center gap-x-4 p-2
                                        cursor-pointer rounded-md transition-colors
                                        ${location.pathname === menu.path ? "bg-white/30" : "hover:bg-white/20"}
                                    `}
                                    onClick={() => {
                                        if (menu.path) navigate(menu.path);
                                        if (menu.subMenu) {
                                            setActiveMenu(activeMenu === index ? null : index);
                                        }
                                    }}
                                >
                                    <span className="text-2xl shrink-0">{menu.icon}</span>
                                    <span className={`flex-1 ${isCollapsed ? "hidden" : ""}`}>{menu.title}</span>

                                    {isCollapsed && (
                                        <span className="absolute left-16 bg-blue-300 text-white text-sm p-2 rounded opacity-0 group-hover:opacity-100 whitespace-nowrap pointer-events-none z-50">
                                            {menu.title}
                                        </span>
                                    )}

                                    {menu.subMenu && !isCollapsed && (
                                        <IoIosArrowDown
                                            className={`duration-300 shrink-0 ${activeMenu === index ? "rotate-180" : ""}`}
                                        />
                                    )}
                                </li>

                                {menu.subMenu && activeMenu === index && !isCollapsed && (
                                    <ul className="pt-1 pb-1">
                                        {menu.subMenuItems.map((sub, i) => (
                                            <li
                                                key={i}
                                                onClick={() => navigate(sub.path)}
                                                className={`
                                                    text-white text-sm p-2 pl-10 cursor-pointer rounded-md mt-1 transition-colors
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

                    <ul>
                        <li
                            className="group relative text-white flex items-center gap-x-4 p-2 cursor-pointer hover:bg-white/50 rounded-md transition-colors"
                            onClick={handleLogout}
                        >
                            <span className="text-2xl shrink-0">{logoutMenu.icon}</span>
                            <span className={`${isCollapsed ? "hidden" : ""}`}>
                                {isPending ? "Logging out..." : logoutMenu.title}
                            </span>
                            {isCollapsed && (
                                <span className="absolute left-16 bg-black text-white text-xs px-2 py-1 rounded opacity-0 group-hover:opacity-100 whitespace-nowrap pointer-events-none z-50">
                                    {isPending ? "Logging out..." : logoutMenu.title}
                                </span>
                            )}
                        </li>
                    </ul>
                </div>
            </aside>

            <main className="flex-1 min-w-0 h-screen overflow-y-auto">
                <div className="p-5 pt-16 lg:pt-8">
                    {children}
                </div>
            </main>
        </div>
    );
}