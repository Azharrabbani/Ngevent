import { useNavigate } from "react-router-dom";

type TabType = "event" | "event_owner";

interface Props {
    activeTab: TabType;
}

export default function EventCreatorTabs({ activeTab }: Props) {
    const navigate = useNavigate();

    return (
        <div
            className="
                flex
                border border-slate-200
                rounded-xl
                p-1
                bg-slate-50
                w-fit
            "
        >
            <button
                onClick={() => navigate("/")}
                className={`
                    px-5 py-2 rounded-lg text-sm
                    transition font-medium
                    ${activeTab === "event"
                        ? "bg-white shadow-sm text-blue-600"
                        : "text-slate-500 hover:text-slate-700"
                    }
                `}
            >
                Event
            </button>

            <button
                onClick={() => navigate("/event-owner")}
                className={`
                    px-5 py-2 rounded-lg text-sm
                    transition font-medium
                    ${activeTab === "event_owner"
                        ? "bg-white shadow-sm text-blue-600"
                        : "text-slate-500 hover:text-slate-700"
                    }
                `}
            >
                Event Owner
            </button>
        </div>
    );
}