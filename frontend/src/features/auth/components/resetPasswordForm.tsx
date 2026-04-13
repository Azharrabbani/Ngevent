import React, { useState } from "react"
import Button from "../../../components/Button"
import PasswordInput from "../../../components/passwordInput"

interface Props {
    onSubmit: (newPassword: string, confirmPassword: string) => void
    loading: boolean
    errors: Record<string, string>
}

export default function ResetPasswordForm({onSubmit, loading, errors}: Props) {
    const [newPassword, setNewPassword] = useState("")
    const [confirmPassword, setConfirmPassword] = useState("")

     const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault()
        onSubmit(newPassword, confirmPassword)
    }

    return (
        <form 
        onSubmit={handleSubmit}
        action="" 
        className="flex flex-col gap-4 mb-2 mt-4">
            <PasswordInput
                type="password"
                name="new_password"
                placeholder="new password"
                onChange={(e) => setNewPassword(e.target.value)}
                error={errors.new_password}
            />
            <PasswordInput
                type="password"
                name="confirm_password"
                placeholder="confirm password"
                onChange={(e) => setConfirmPassword(e.target.value)}
                error={errors.confirm_password}
            />
            <Button 
            type="submit">
                {loading? "loading..." : "Submit" }
            </Button>
        </form>
    )
}