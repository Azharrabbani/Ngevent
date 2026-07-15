import { useEffect, useRef, useState } from "react";
import { CrossIcon } from "../../../components/icon";
import Slider from "./slider";
import Input from "../../../components/input";
import Button from "../../../components/Button";
import PasswordInput from "../../../components/passwordInput";
import { useRegisterAdmin } from "../../auth/hooks/useRegisterAdmin";

interface Props {
    isOpen: boolean;
    onClose: () => void;
}

export default function CreateAdminSlideOver({ isOpen, onClose }: Props) {
    const inputRef = useRef<HTMLInputElement>(null);
    const [email, setEmail] = useState<string>("");
    const [password, setPassword] = useState<string>("");
    const [confirmPassword, setConfirmPassword] = useState<string>("");
    const [showPassword, setShowPassword] = useState(false);
    const [showConfirmPassword, setShowConfirmPassword] = useState(false);

    const {
        mutateAsync: register,
        isPending,
        error,
    } = useRegisterAdmin();

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        if (!email.trim() && !password.trim() && !confirmPassword.trim()) return;
        register({ email: email, password: password, confirm_password: confirmPassword });
    };

    const validationErrors =
        (error as any)?.response?.data?.errors || {};

    useEffect(() => {
        if (isOpen) {
            setTimeout(() => inputRef.current?.focus(), 150);
        } else {
            setEmail("");
            setPassword("");
            setConfirmPassword("");

            setShowPassword(false);
            setShowConfirmPassword(false);
        }
    }, [isOpen]);

    useEffect(() => {
        const handleKey = (e: KeyboardEvent) => {
            if (e.key === "Escape" && isOpen) onClose();
        };
        document.addEventListener("keydown", handleKey);
        return () => document.removeEventListener("keydown", handleKey);
    }, [isOpen, onClose]);

    return (
        <Slider isOpen={isOpen} onClose={onClose}>
            <>
                <div className="flex items-center justify-between px-6 py-5 border-b border-gray-100">
                    <div>
                        <h2 className="text-xl font-bold text-gray-900">New Admin</h2>
                        <p className="text-sm text-gray-500 mt-0.5">
                            Add a new admin to the system
                        </p>
                    </div>
                    <button
                        onClick={onClose}
                        className="w-9 h-9 flex items-center justify-center rounded-xl hover:bg-gray-100 transition-colors text-gray-500 hover:text-gray-800"
                    >
                        <CrossIcon className="w-5 h-5" />
                    </button>
                </div>

                <form
                    onSubmit={handleSubmit}
                    className="flex flex-col flex-1 px-6 py-6 gap-6">
                    <div className="flex flex-col gap-1.5">
                        <label className="text-sm font-semibold text-gray-700">
                            Email <span className="text-red-500">*</span>
                        </label>
                        <Input
                            ref={inputRef}
                            type="email"
                            placeholder="e.g. [EMAIL_ADDRESS]"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            error={validationErrors.email}
                            className="w-full bg-white px-4 py-3 rounded-lg border border-gray-300 outline-none focus:ring-2 focus:ring-blue-500 text-sm"
                        />
                    </div>

                    <div className="flex flex-col gap-1.5">
                        <label className="text-sm font-semibold text-gray-700">
                            Password <span className="text-red-500">*</span>
                        </label>
                        <PasswordInput
                            show={showPassword}
                            onShowChange={setShowPassword}
                            placeholder="password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            error={validationErrors.password}
                            className="w-full bg-white px-4 py-3 rounded-lg border border-gray-300 outline-none focus:ring-2 focus:ring-blue-500 text-sm"
                        />
                    </div>

                    <div className="flex flex-col gap-1.5">
                        <label className="text-sm font-semibold text-gray-700">
                            Confirm Password <span className="text-red-500">*</span>
                        </label>
                        <PasswordInput
                            show={showConfirmPassword}
                            onShowChange={setShowConfirmPassword}
                            placeholder="confirm password"
                            value={confirmPassword}
                            onChange={(e) => setConfirmPassword(e.target.value)}
                            error={validationErrors.confirm_password}
                            className="w-full bg-white px-4 py-3 rounded-lg border border-gray-300 outline-none focus:ring-2 focus:ring-blue-500 text-sm"
                        />
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
                            disabled={
                                !email.trim() || !password.trim() || !confirmPassword.trim() || isPending
                            }
                            className="flex-1 px-4 py-3 rounded-xl bg-[#0066FF] text-white font-semibold text-sm disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                        >
                            {isPending ? (
                                <>
                                    <span className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                                    Registering...
                                </>
                            ) : (
                                "Register Admin"
                            )}
                        </Button>
                    </div>
                </form>
            </>
        </Slider>
    )
}