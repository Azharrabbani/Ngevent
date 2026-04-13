import type React from "react"
import { cn } from "../../../utils/cn"

interface AuthContainerProps  {
    children: React.ReactNode
    className?: string
}

export default function AuthContainer({children, className=""}: AuthContainerProps) {
    return (
        <section className="bg-gray-50 min-h-screen flex items-center justify-center">
            <div className={cn(
                "bg-gray-100 flex rounded-2xl shadow-lg  max-w-4xl items-center",
                className
            )}>
                {children}
            </div>
        </section>
    )
}