type EventStatus = "active" | "done";

interface Props {
    activeStatus: EventStatus;
    onChange: (status: EventStatus) => void;
}

export default function EventStatusTabs({ activeStatus, onChange }: Props) {
    return (
        <div className="flex border border-slate-200 rounded-xl p-1 bg-slate-50 w-fit">
            <button
                onClick={() => onChange("active")}
                className={`
                    px-5 py-2 rounded-lg text-sm transition font-medium
                    ${activeStatus === "active"
                        ? "bg-slate-800 text-white shadow-sm"
                        : "text-slate-500 hover:text-slate-700"
                    }
                `}
            >
                Active events
            </button>
            <button
                onClick={() => onChange("done")}
                className={`
                    px-5 py-2 rounded-lg text-sm transition font-medium
                    ${activeStatus === "done"
                        ? "bg-slate-800 text-white shadow-sm"
                        : "text-slate-500 hover:text-slate-700"
                    }
                `}
            >
                Past events
            </button>
        </div>
    );
}