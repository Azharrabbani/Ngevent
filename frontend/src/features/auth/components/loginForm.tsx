import { useState } from "react";
import Button from "../../../components/Button";
import Input from "../../../components/input";
import PasswordInput from "../../../components/passwordInput";

interface Props {
    onSubmit: (email: string, password: string) => void
    loading: boolean
    errors: Record<string, string>
}

export default function LoginForm({ onSubmit, loading, errors }: Props) {
    const [email, setEmail] = useState<string>("");
    const [password, setPassword] = useState<string>("");

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault()
        onSubmit(email, password)
    };

    return (
        <form onSubmit={handleSubmit} className="flex flex-col gap-5 mb-2">
            <div className="flex flex-col gap-1.5">
                <label htmlFor="email" className="text-sm font-medium text-gray-600">
                    Email
                </label>
                <Input
                    id="email"
                    type="email"
                    name="email"
                    placeholder="your-email@email.com"
                    className="h-11 w-96"
                    onChange={(e) => setEmail(e.target.value)}
                    error={errors.email}
                />
            </div>

            <div className="flex flex-col gap-1.5">
                <label htmlFor="password" className="text-sm font-medium text-gray-600">
                    Password
                </label>
                <PasswordInput
                    id="password"
                    name="password"
                    placeholder="password"
                    className="h-11 w-96"
                    onChange={(e) => setPassword(e.target.value)}
                    error={errors.password}
                />
            </div>

            <Button disabled={loading} type="submit" className="h-11 mt-1">
                {loading ? "Loading..." : "Login"}
            </Button>
        </form>
    )
}