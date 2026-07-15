import { useState, useEffect } from "react";
import type { EventCategory } from "../../../types/publicEventResponse";

interface Props {
    banner: string;
    name: string;
    categories: EventCategory[];
}

export default function EventBanner({ banner, name, categories }: Props) {
    const [isPreviewOpen, setIsPreviewOpen] = useState(false);

    // Lock body scroll saat preview terbuka
    useEffect(() => {
        if (isPreviewOpen) {
            document.body.style.overflow = "hidden";
        } else {
            document.body.style.overflow = "";
        }
        return () => {
            document.body.style.overflow = "";
        };
    }, [isPreviewOpen]);

    // Tutup preview dengan tombol Escape
    useEffect(() => {
        if (!isPreviewOpen) return;

        const handleKeyDown = (e: KeyboardEvent) => {
            if (e.key === "Escape") setIsPreviewOpen(false);
        };

        window.addEventListener("keydown", handleKeyDown);
        return () => window.removeEventListener("keydown", handleKeyDown);
    }, [isPreviewOpen]);

    return (
        <>
            <div className="relative rounded-2xl overflow-hidden bg-slate-100 aspect-[16/9] w-full">
                <button
                    type="button"
                    onClick={() => setIsPreviewOpen(true)}
                    className="w-full h-full cursor-zoom-in group"
                    aria-label="Preview banner image"
                >
                    <img
                        src={banner}
                        alt={name}
                        className="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
                    />
                </button>
                <div className="absolute inset-0 bg-gradient-to-t from-black/20 to-transparent pointer-events-none" />

                {categories.length > 0 && (
                    <div className="absolute top-3 left-3 flex flex-wrap gap-2 pointer-events-none">
                        {categories.map((cat) => (
                            <span
                                key={cat.id}
                                className="inline-flex items-center gap-1.5 bg-white/90 backdrop-blur-sm text-slate-700 text-xs font-medium px-3 py-1.5 rounded-full shadow-sm"
                            >
                                <span className="text-[10px]">✦</span>
                                {cat.name}
                            </span>
                        ))}
                    </div>
                )}
            </div>

            {isPreviewOpen && (
                <div
                    className="fixed inset-0 z-50 bg-black/90 flex items-center justify-center p-4 cursor-zoom-out"
                    onClick={() => setIsPreviewOpen(false)}
                >
                    <button
                        type="button"
                        onClick={() => setIsPreviewOpen(false)}
                        className="absolute top-4 right-4 text-white/80 hover:text-white bg-white/10 hover:bg-white/20 rounded-full w-10 h-10 flex items-center justify-center transition-colors"
                        aria-label="Close preview"
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            className="h-5 w-5"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                            strokeWidth={2}
                        >
                            <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>

                    <img
                        src={banner}
                        alt={name}
                        className="max-w-full max-h-full object-contain rounded-lg"
                        onClick={(e) => e.stopPropagation()}
                    />
                </div>
            )}
        </>
    );
}