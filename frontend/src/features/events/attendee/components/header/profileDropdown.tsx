import { useState, useRef, useEffect } from "react";
import { FaChevronDown } from "react-icons/fa";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../../../../../lib/auth";
import { useLogout } from "../../../../auth/hooks/useLogout";
import { UserIcon } from "../../../../../components/icon";

export default function ProfileDropdown() {
    const navigate = useNavigate();
    const { user } = useAuth();
    const logoutMutation = useLogout();

    const [open, setOpen] = useState(false);
    const ref = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const handleClickOutside = (e: MouseEvent) => {
            if (
                ref.current &&
                !ref.current.contains(e.target as Node)
            ) {
                setOpen(false);
            }
        };

        document.addEventListener("mousedown", handleClickOutside);
        return () => document.removeEventListener("mousedown", handleClickOutside);
    }, []);

    const displayName = user?.role === "admin" ? "Admin" : "Event Owner";

    return (
        <div ref={ref} className="relative">
            <button
                onClick={() => setOpen(!open)}
                className="
                    flex items-center gap-3
                    border border-slate-200 rounded-full
                    px-4 py-2 bg-slate-50/50
                    hover:bg-slate-50 transition duration-150
                "
            >
                <UserIcon className="w-5 h-8" />
                <span className="text-sm font-medium text-slate-700 capitalize">
                    {displayName}
                </span>

                <FaChevronDown className="w-3.5 h-3.5 text-slate-500" />
            </button>

            {open && (
                <div
                    className="
                        absolute right-0 mt-2
                        bg-white border border-slate-200
                        rounded-xl shadow-lg
                        w-52 overflow-hidden
                        z-50
                    "
                >
                    {user?.role === "admin" && (
                        <button
                            onClick={() => {
                                setOpen(false);
                                navigate("/admin/dashboard");
                            }}
                            className="
                                w-full text-left
                                px-4 py-3 text-sm text-slate-700
                                hover:bg-slate-50 hover:text-blue-600
                                border-b border-slate-100 transition
                            "
                        >
                            Admin Dashboard
                        </button>
                    )}

                    {user?.role === "event organizer" && (
                        <div>
                            <button
                                onClick={() => {
                                    setOpen(false);
                                    navigate("/organizer/dashboard");
                                }}
                                className="
                                    w-full text-left
                                    px-4 py-3 text-sm text-slate-700
                                    hover:bg-slate-50 hover:text-blue-600
                                    border-b border-slate-100 transition
                                "
                            >
                                Organizer Dashboard
                            </button>

                            <button
                                onClick={() => {
                                    setOpen(false);
                                    navigate("/profile");
                                }}
                                className="
                            w-full text-left
                            px-4 py-3 text-sm   text-slate-700
                            hover:bg-slate-50 hover:text-blue-600
                            border-b border-slate-100 transition
                        "
                            >
                                My Profile
                            </button>
                        </div>

                    )}

                    <button
                        onClick={() => {
                            setOpen(false);
                            logoutMutation.mutate();
                        }}
                        className="
                            w-full text-left
                            px-4 py-3 text-sm text-red-600
                            hover:bg-red-50 transition font-medium
                        "
                    >
                        Logout
                    </button>
                </div>
            )}
        </div>
    );
}