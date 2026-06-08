import { useState } from "react";

type TabType =
    | "event"
    | "event_creator";

export default function EventCreatorTabs() {
    const [activeTab, setActiveTab] =
        useState<TabType>("event");

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
                onClick={() =>
                    setActiveTab("event")
                }
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
                onClick={() =>
                    setActiveTab(
                        "event_creator"
                    )
                }
                className={`
                    px-5 py-2 rounded-lg text-sm
                    transition font-medium
                    ${activeTab ===
                        "event_creator"
                        ? "bg-white shadow-sm text-blue-600"
                        : "text-slate-500 hover:text-slate-700"
                    }
                `}
            >
                Event Creator
            </button>
        </div>
    );
}