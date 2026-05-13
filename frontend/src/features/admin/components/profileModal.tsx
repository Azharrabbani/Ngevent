interface Props {
    isOpen: boolean;
    onClose: () => void;
    children: React.ReactElement
    isLoading: boolean;
}

export default function ProfileModal({
    isOpen,
    onClose,
    children,
    isLoading,
}: Props) {
    if (!isOpen) return null;

    return (
        <div
            className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4"
            onClick={onClose}
        >
            <div
                onClick={(e) => e.stopPropagation()}
                className="
                    w-full max-w-4xl
                    max-h-[90vh] overflow-y-auto
                    bg-white border-2 border-black rounded-xl
                    p-4 sm:p-6
                    shadow-[6px_6px_0px_black]
                    relative
                "
            >
                <button
                    onClick={onClose}
                    className="absolute top-2 right-2 border-2 border-black px-2 rounded hover:bg-black hover:text-white"
                >
                    ✕
                </button>

                {isLoading && <h1 className="text-center">Loading...</h1>}
                
                {children}
            </div>
        </div>
    );
}