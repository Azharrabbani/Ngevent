import toast from "react-hot-toast";
import Button from "./Button";

export default function ShowToast(message: string) {
    return (
        toast((t) => (
            <span>
                {message}
                <Button
                    onClick={() => toast.dismiss(t.id)}
                    className="text-red-500 font-bold bg-white"
                >
                    ✕
                </Button>
            </span>
            
        ))
    )
}