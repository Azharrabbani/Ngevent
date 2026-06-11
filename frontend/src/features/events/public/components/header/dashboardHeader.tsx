import { Link, useNavigate } from "react-router-dom";
import { FiSearch } from "react-icons/fi";
import ProfileDropdown from "./profileDropdown";
import { useAuth } from "../../../../../lib/auth";
import { useEffect, useState } from "react";

interface Props {
    onSearchChange?: (val: string) => void;
    isEventOwnerDashboard?: boolean;
}

export default function DashboardHeader({
    onSearchChange,
    isEventOwnerDashboard = false,
}: Props) {
    const navigate = useNavigate();
    const { user, loading } = useAuth();
    const [inputValue, setInputValue] = useState("");

    useEffect(() => {
        if (!onSearchChange) return;
        const handler = setTimeout(() => {
            onSearchChange(inputValue);
        }, 450);

        return () => {
            clearTimeout(handler);
        };
    }, [inputValue, onSearchChange]);

    return (
        <header
            className="
                sticky top-0 z-40
                bg-white/95 backdrop-blur-sm
                border-b border-slate-200
                shadow-sm
            "
        >
            <div className="max-w-7xl mx-auto px-4 lg:px-6 h-16 flex items-center justify-between gap-6">
                <Link
                    to="/"
                    className="
                        text-2xl font-extrabold
                        bg-gradient-to-r from-blue-600 to-indigo-600
                        bg-clip-text text-transparent
                        shrink-0
                        tracking-tight
                    "
                >
                    Ngevent
                </Link>

                <div className="flex-1 max-w-lg relative">
                    <FiSearch className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400 w-4 h-4" />
                    {isEventOwnerDashboard ? (
                        <input
                            type="text"
                            value={inputValue}
                            onChange={(e) => setInputValue(e.target.value)}
                            placeholder="Find your favorite events..."
                            className="
                                w-full pl-9 pr-4 py-2.5
                                text-sm
                                bg-slate-50 border border-slate-200
                                rounded-xl
                                focus:outline-none focus:ring-2 focus:ring-blue-500/30 focus:border-blue-400
                                placeholder:text-slate-400
                                transition
                            "
                        />

                    ) : (
                        <input
                            type="text"
                            value={inputValue}
                            onChange={(e) => setInputValue(e.target.value)}
                            placeholder="Find your favorite event owners..."
                            className="
                                w-full pl-9 pr-4 py-2.5
                                text-sm
                                bg-slate-50 border border-slate-200
                                rounded-xl
                                focus:outline-none focus:ring-2 focus:ring-blue-500/30 focus:border-blue-400
                                placeholder:text-slate-400
                                transition
                            "
                        />
                    )}
                </div>

                {loading ? (
                    <div className="w-24 h-9 bg-slate-100 animate-pulse rounded-lg shrink-0" />
                ) : user ? (
                    <ProfileDropdown />
                ) : (
                    <div className="flex items-center gap-2 shrink-0">
                        <button
                            onClick={() => navigate("/login")}
                            className="
                                px-4 py-2 text-sm font-medium
                                text-slate-700 hover:text-blue-600
                                hover:bg-slate-100
                                rounded-lg transition
                            "
                        >
                            Login
                        </button>
                        <button
                            onClick={() => navigate("/register")}
                            className="
                                px-4 py-2 text-sm font-semibold
                                text-white
                                bg-gradient-to-r from-blue-600 to-indigo-600
                                hover:from-blue-700 hover:to-indigo-700
                                rounded-lg transition shadow-sm
                            "
                        >
                            Register
                        </button>
                    </div>
                )}
            </div>
        </header>
    );
}