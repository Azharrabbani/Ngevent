import type React from "react"
import { cn } from "../../../utils/cn"

interface CompleteProfileProps {
    children: React.ReactNode
    className?: string
};

export default function CompleteProfileContainer({children, className=""}: CompleteProfileProps) {
    return (            
        <section className="min-h-screen flex items-center justify-center">
            <div className={cn(
                "flex max-w-4xl items-center mx-auto p-20",
                className
            )}>
                {children}
            </div>
        </section>
    )
}