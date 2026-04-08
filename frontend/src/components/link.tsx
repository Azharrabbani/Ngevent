import type React from "react"

interface LinkProps  {
    children: React.ReactNode
    endpoint: string
    className?: string
}

export default function Link(
    {
        children, 
        className="", 
        endpoint
    }: LinkProps) {
    return (
        <a 
        className={`text-blue-500 hover:text-blue-700 text-sm ${className}`}
        href={`${endpoint}`}>
            {children}
        </a>
    )
}