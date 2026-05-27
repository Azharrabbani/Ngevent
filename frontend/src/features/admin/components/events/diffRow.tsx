import React from "react";

interface DiffRowProps {
    label: string;
    currentContent: React.ReactNode;
    proposedContent: React.ReactNode;
    hasChange?: boolean;
}

export default function DiffRow({
    label,
    currentContent,
    proposedContent,
    hasChange = false,
}: DiffRowProps) {
    return (
        <div className="grid grid-cols-1 sm:grid-cols-2 border-b border-gray-100 last:border-b-0">
            {/* Current */}
            <div className="p-4 sm:border-r border-gray-100 border-b sm:border-b-0">
                <span className="block text-[10px] font-bold uppercase tracking-widest text-gray-400 mb-1.5">
                    {label}
                </span>
                <div className="text-sm text-gray-700">{currentContent}</div>
            </div>

            {/* Proposed */}
            <div
                className={`p-4 ${hasChange ? "bg-amber-50/60" : ""
                    }`}
            >
                <span className="block text-[10px] font-bold uppercase tracking-widest text-gray-400 mb-1.5">
                    {label}
                </span>
                <div className={`text-sm ${hasChange ? "text-gray-900 font-medium" : "text-gray-700"}`}>
                    {proposedContent}
                </div>
                {hasChange && (
                    <span className="inline-block mt-1.5 text-[10px] px-1.5 py-0.5 rounded bg-amber-100 text-amber-700 font-semibold uppercase tracking-wide">
                        Changed
                    </span>
                )}
            </div>
        </div>
    );
}