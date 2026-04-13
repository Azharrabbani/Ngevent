import { useState } from "react";
import Button from "../../../components/Button";
import Input from "../../../components/input";

interface Props {
    onSubmit: (email: string) => void
    loading: boolean
    errors: Record<string, string>
}

export default function ForgetPasswordForm({onSubmit, loading, errors}: Props) {
    const [email, setEmail] = useState("")

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault()
        onSubmit(email)
    }
    
    return (
        <form 
        onSubmit={handleSubmit}
        action="" 
        className="flex flex-col gap-4 mb-2">
        <Input
            type="email"
            name="email"
            placeholder="email"
            className="mt-8"
            onChange={(e) => setEmail(e.target.value)}
            error={errors.email}
        />
        <Button 
        disabled={loading}
        type="submit">
            {loading ? "Loading..." : "Submit"}
        </Button>
        </form>
    )
}