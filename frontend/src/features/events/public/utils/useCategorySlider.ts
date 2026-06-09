import { useEffect, useState } from "react";

export default function useCategorySlider(
    ref: React.RefObject<HTMLDivElement | null>
) {
    const [canScrollLeft, setCanScrollLeft] =
        useState(false);

    const [canScrollRight, setCanScrollRight] =
        useState(false);

    const checkScroll = () => {
        if (!ref.current) return;

        const el = ref.current;

        setCanScrollLeft(el.scrollLeft > 0);

        setCanScrollRight(
            el.scrollLeft <
            el.scrollWidth - el.clientWidth - 5
        );
    };

    const scrollLeft = () => {
        ref.current?.scrollBy({
            left: -300,
            behavior: "smooth",
        });
    };

    const scrollRight = () => {
        ref.current?.scrollBy({
            left: 300,
            behavior: "smooth",
        });
    };

    useEffect(() => {
        checkScroll();

        const el = ref.current;

        if (!el) return;

        el.addEventListener(
            "scroll",
            checkScroll
        );

        window.addEventListener(
            "resize",
            checkScroll
        );

        return () => {
            el.removeEventListener(
                "scroll",
                checkScroll
            );

            window.removeEventListener(
                "resize",
                checkScroll
            );
        };
    }, [ref]);

    return {
        canScrollLeft,
        canScrollRight,
        scrollLeft,
        scrollRight,
    };
}