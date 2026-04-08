import type React from "react"

interface AuthContainerProps  {
    children: React.ReactNode
}

export default function AuthContainer({children}: AuthContainerProps) {
    return (
        <section className="bg-gray-50 min-h-screen flex items-center justify-center">
            <div className="bg-gray-100 flex rounded-2xl shadow-lg  max-w-4xl items-center">
                {children}
            </div>
        </section>
    )
}