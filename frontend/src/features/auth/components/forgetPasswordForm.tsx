import Button from "../../../components/Button";
import Input from "../../../components/input";

export default function ForgetPasswordForm() {
    return (
        <form action="" 
        className="flex flex-col gap-4 mb-2">
        <Input
            type="email"
            name="email"
            placeholder="email"
            className="mt-8"
        />
        <Button type="submit">
            Send email
        </Button>
        </form>
    )
}