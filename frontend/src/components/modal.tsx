interface Props {
    children: React.ReactNode;
    isOpen: boolean;
    onClose: () => void;
};

export default function Modal({ children, isOpen, onClose }: Props) {
    if (!isOpen) return null;

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
            <div
                className="absolute inset-0 bg-black/50 backdrop-blur-sm"
                onClick={onClose}
            />

            {children}
        </div>
    )
}