interface Props {
    isOpen: boolean;
    onClose: () => void;
    children: React.ReactElement;
};

export default function Slider({ isOpen, onClose, children }: Props) {
    return (
        <>
            <div
                className={`fixed inset-0 z-40 bg-black/30 backdrop-blur-[2px] transition-opacity duration-300 ${isOpen ? "opacity-100 pointer-events-auto" : "opacity-0 pointer-events-none"
                    }`}
                onClick={onClose}
            />

            <div className={`fixed top-0 right-0 z-50 h-full w-full max-w-md bg-white shadow-2xl flex flex-col transition-transform duration-300 ease-in-out ${isOpen ? "translate-x-0" : "translate-x-full"
                }`}>
                {children}
            </div>
        </>
    );
};