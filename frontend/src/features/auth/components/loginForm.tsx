import { useState } from "react";
import Button from "../../../components/Button";
import Input from "../../../components/input";
import PasswordInput from "../../../components/passwordInput";

interface Props {
    onSubmit: (email: string, password: string) => void
    loading: boolean
    errors: Record<string, string>
}

export default function LoginForm({onSubmit, loading, errors}: Props) {
    const [email, setEmail] = useState<string>("");
    const [password, setPassword] = useState<string>("");

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault()
        onSubmit(email, password)
    };

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
                name="password"
                placeholder="password"
                onChange={(e) => setPassword(e.target.value)}
                error={errors.password}
            />
            <Button
                disabled={loading}  
                type="submit"
            >
                {loading ? "Loading..." : "Login"}    
            </Button>
        </form>
    )
}