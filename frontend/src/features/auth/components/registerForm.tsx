import React, { useState } from "react";
import Button from "../../../components/Button";
import Input from "../../../components/input";
import PasswordInput from "../../../components/passwordInput";

interface Props {
    onSubmit: (email: string, password: string, conFirmPassword: string) => void
    loading: boolean
    errors: Record<string, string>
}

export default function RegisterForm({onSubmit, loading, errors}: Props) {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [confirmPassword, setConfirmPassword] = useState("");

    const handleSubmit = ((e: React.FormEvent) => {
        e.preventDefault();
        onSubmit(email, password, confirmPassword);
    })

    return (
        <form 
        onSubmit={handleSubmit} 
        className="flex flex-col gap-4 mb-2">
            <Input
                type="email"
                name="email"
                placeholder="email"
                className="mt-8"
                onChange={(e) => setEmail(e.target.value)}
                error={errors.email}
            />
            <PasswordInput
                type="password"
                name="password"
                placeholder="password"
                onChange={(e) => setPassword(e.target.value)}
                error={errors.password}
            />
            <PasswordInput
                type="password"
                name="confirm_password"
                placeholder="comfirm password"
                onChange={(e) => setConfirmPassword(e.target.value)}
                error={errors.confirm_password}
            />
            <Button 
            disabled={loading}
            type="submit">
                {loading ? "loading..." : "Register"}
            </Button>
        </form>
    )
}