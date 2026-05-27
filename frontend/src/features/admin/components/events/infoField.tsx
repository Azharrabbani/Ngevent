import React from "react";

interface InfoFieldProps {
    label: string;
    children: React.ReactNode;
    className?: string;
}

export default function InfoField({ label, children, className = "" }: InfoFieldProps) {
    return (
        <div className={`flex flex-col gap-0.5 ${className}`}>
            <span className="text-[11px] font-semibold uppercase tracking-widest text-gray-400">
                {label}
            </span>
            <div className="text-sm text-gray-800 font-medium">{children}</div>
        </div>
    )
}