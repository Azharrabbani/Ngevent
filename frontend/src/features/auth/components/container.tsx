import { cn } from "../../../utils/cn";

interface AuthContainerProps {
    children: React.ReactNode;
    className?: string;
}

export default function AuthContainer({
    children,
    className = "",
}: AuthContainerProps) {
    return (
        <section
            className="min-h-screen flex items-center justify-center px-4"
            style={{
                backgroundImage: "url('/auth-bg.webp')",
                backgroundSize: "cover",
                backgroundPosition: "center",
            }}
        >
            <div
                className={cn(
                    "w-full max-w-md bg-white/90 backdrop-blur-lg rounded-3xl shadow-xl border border-white/30",
                    className
                )}
            >
                {children}
            </div>
        </section>
    );
}