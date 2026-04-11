import React, { useState } from "react";
import Button from "../../../components/Button";
import Input from "../../../components/input";

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
            <Input
                type="password"
                name="password"
                placeholder="password"
                rightIcon={
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="gray" className="bi bi-eye absolute top-1/2 right-3 -translate-y-1/2" viewBox="0 0 16 16">
                        <path d="M16 8s-3-5.5-8-5.5S0 8 0 8s3 5.5 8 5.5S16 8 16 8M1.173 8a13 13 0 0 1 1.66-2.043C4.12 4.668 5.88 3.5 8 3.5s3.879 1.168 5.168 2.457A13 13 0 0 1 14.828 8q-.086.13-.195.288c-.335.48-.83 1.12-1.465 1.755C11.879 11.332 10.119 12.5 8 12.5s-3.879-1.168-5.168-2.457A13 13 0 0 1 1.172 8z"/>
                        <path d="M8 5.5a2.5 2.5 0 1 0 0 5 2.5 2.5 0 0 0 0-5M4.5 8a3.5 3.5 0 1 1 7 0 3.5 3.5 0 0 1-7 0"/>
                    </svg>
                }
                onChange={(e) => setPassword(e.target.value)}
                error={errors.password}
            />
            <Input
                type="password"
                name="confirm_password"
                placeholder="comfirm password"
                rightIcon={
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="gray" className="bi bi-eye absolute top-1/2 right-3 -translate-y-1/2" viewBox="0 0 16 16">
                        <path d="M16 8s-3-5.5-8-5.5S0 8 0 8s3 5.5 8 5.5S16 8 16 8M1.173 8a13 13 0 0 1 1.66-2.043C4.12 4.668 5.88 3.5 8 3.5s3.879 1.168 5.168 2.457A13 13 0 0 1 14.828 8q-.086.13-.195.288c-.335.48-.83 1.12-1.465 1.755C11.879 11.332 10.119 12.5 8 12.5s-3.879-1.168-5.168-2.457A13 13 0 0 1 1.172 8z"/>
                        <path d="M8 5.5a2.5 2.5 0 1 0 0 5 2.5 2.5 0 0 0 0-5M4.5 8a3.5 3.5 0 1 1 7 0 3.5 3.5 0 0 1-7 0"/>
                    </svg>
                }
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