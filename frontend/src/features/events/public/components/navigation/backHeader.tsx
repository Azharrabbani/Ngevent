import { FiArrowLeft } from "react-icons/fi";

interface Props {
    onBack: () => void;
}

export default function BackHeader({ onBack }: Props) {
    return (
        <header className="bg-white border-b border-slate-200 sticky top-0 z-10">
            <div className="max-w-6xl mx-auto px-4 lg:px-6 h-14 flex items-center gap-4">
                <button
                    onClick={onBack}
                    className="flex items-center gap-1.5 text-sm text-slate-600 hover:text-slate-900 transition-colors"
                >
                    <FiArrowLeft className="w-4 h-4" />
                    Back
                </button>
            </div>
        </header>
    );
}