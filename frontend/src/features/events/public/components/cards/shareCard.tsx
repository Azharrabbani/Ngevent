import { useState } from "react";
import { FiShare2, FiLink, FiMail } from "react-icons/fi";

interface Props {
    eventName: string;
    eventUrl?: string;
}

export default function ShareCard({ eventName, eventUrl }: Props) {
    const [copied, setCopied] = useState(false);

    const url = eventUrl ?? window.location.href;

    const handleNativeShare = async () => {
        if (navigator.share) {
            try {
                await navigator.share({ title: eventName, url });
            } catch {
            }
        }
    };

    const handleCopyLink = async () => {
        try {
            await navigator.clipboard.writeText(url);
            setCopied(true);
            setTimeout(() => setCopied(false), 2000);
        } catch {
        }
    };

    const handleEmail = () => {
        const subject = encodeURIComponent(`Check out: ${eventName}`);
        const body = encodeURIComponent(`I thought you'd be interested in this event: ${url}`);
        window.open(`mailto:?subject=${subject}&body=${body}`);
    };

    const actions = [
        {
            icon: <FiShare2 className="w-4 h-4" />,
            label: "Share",
            onClick: handleNativeShare,
            title: "Share via apps",
        },
        {
            icon: <FiLink className="w-4 h-4" />,
            label: copied ? "Copied!" : "Copy",
            onClick: handleCopyLink,
            title: "Copy link",
        },
        {
            icon: <FiMail className="w-4 h-4" />,
            label: "Email",
            onClick: handleEmail,
            title: "Share via email",
        },
    ];

    return (
        <div className="bg-white rounded-2xl border border-slate-200 p-6 space-y-4">
            <p className="text-xs font-semibold uppercase tracking-widest text-slate-400">
                Share with friends
            </p>
            <div className="flex items-center gap-3">
                {actions.map((action) => (
                    <button
                        key={action.label}
                        onClick={action.onClick}
                        title={action.title}
                        className="flex flex-col items-center gap-1.5 group"
                    >
                        <span className="w-10 h-10 rounded-xl border border-slate-200 flex items-center justify-center text-slate-600 group-hover:bg-blue-50 group-hover:border-blue-200 group-hover:text-blue-600 transition-colors">
                            {action.icon}
                        </span>
                        <span className="text-[10px] text-slate-400 group-hover:text-blue-500 transition-colors">
                            {action.label}
                        </span>
                    </button>
                ))}
            </div>
        </div>
    );
}