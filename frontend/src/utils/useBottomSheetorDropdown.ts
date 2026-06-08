import { useState, useRef, useEffect, useCallback } from "react";

export function useBottomSheetOrDropdown() {
    const [open, setOpen] = useState(false);
    const [dropdownPos, setDropdownPos] = useState({ top: 0, left: 0 });
    const buttonRef = useRef<HTMLButtonElement>(null);
    const dropdownRef = useRef<HTMLDivElement>(null);
    const [isMobile, setIsMobile] = useState(() => window.innerWidth < 1024);

    useEffect(() => {
        const handler = () => setIsMobile(window.innerWidth < 1024);
        window.addEventListener("resize", handler);
        return () => window.removeEventListener("resize", handler);
    }, []);

    const updatePosition = useCallback(() => {
        if (!buttonRef.current) return;
        const rect = buttonRef.current.getBoundingClientRect();
        setDropdownPos({
            top: rect.bottom + window.scrollY + 4,
            left: Math.min(
                rect.right + window.scrollX - 220,
                window.innerWidth - 232
            ),
        });
    }, []);

    const handleToggle = () => {
        if (!open && !isMobile) updatePosition();
        setOpen((prev) => !prev);
    };

    const close = () => setOpen(false);

    useEffect(() => {
        if (!open || isMobile) return;
        const handler = (e: MouseEvent) => {
            if (
                dropdownRef.current?.contains(e.target as Node) ||
                buttonRef.current?.contains(e.target as Node)
            ) return;
            close();
        };
        document.addEventListener("mousedown", handler);
        return () => document.removeEventListener("mousedown", handler);
    }, [open, isMobile]);

    useEffect(() => {
        if (!open || isMobile) return;
        window.addEventListener("scroll", updatePosition, true);
        window.addEventListener("resize", updatePosition);
        return () => {
            window.removeEventListener("scroll", updatePosition, true);
            window.removeEventListener("resize", updatePosition);
        };
    }, [open, isMobile, updatePosition]);

    useEffect(() => {
        document.body.style.overflow = isMobile && open ? "hidden" : "";
        return () => { document.body.style.overflow = ""; };
    }, [isMobile, open]);

    return { open, isMobile, handleToggle, close, buttonRef, dropdownRef, dropdownPos };
}