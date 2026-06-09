import { FiCheckCircle } from "react-icons/fi";

interface Props {
    name: string;
    photoProfile: string;
    isVerified: boolean;
    onViewProfile?: () => void;
}

export default function EventOwnerCard({ name, photoProfile, isVerified, onViewProfile }: Props) {
    return (
        <div className="bg-white rounded-2xl border border-slate-200 p-6 space-y-4">
            <p className="text-xs font-semibold uppercase tracking-widest text-slate-400">
                Presented by
            </p>

            <div className="flex items-center gap-3">
                <img
                    src={photoProfile}
                    alt={name}
                    className="w-12 h-12 rounded-full object-cover ring-2 ring-slate-100"
                    onError={(e) => {
                        (e.currentTarget as HTMLImageElement).src =
                            `https://ui-avatars.com/api/?name=${encodeURIComponent(name)}&background=e0e7ff&color=4f46e5&size=96`;
                    }}
                />
                <div>
                    <div className="flex items-center gap-1.5">
                        <span className="text-sm font-semibold text-slate-900">{name}</span>
                        {isVerified && (
                            <FiCheckCircle className="w-4 h-4 text-blue-500 shrink-0" />
                        )}
                    </div>
                    <p className="text-xs text-slate-500 mt-0.5">Event Owner</p>
                </div>
            </div>

            <button
                onClick={onViewProfile}
                className="w-full rounded-xl border border-slate-200 text-sm font-medium text-slate-700 py-2.5 hover:bg-slate-50 hover:border-slate-300 transition-colors"
            >
                View Event Owner
            </button>
        </div>
    );
}