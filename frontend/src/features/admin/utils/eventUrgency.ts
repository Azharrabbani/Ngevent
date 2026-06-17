import type { EventsResponse } from "../../events/types/eventResponse";

export type UrgencyLevel = "critical" | "urgent" | "warning" | "normal";

export interface UrgencyMeta {
    level: UrgencyLevel;
    daysUntil: number;
    label: string;
}

export interface UrgencyGroup {
    level: UrgencyLevel;
    title: string;
    color: string;
    badgeColor: string;
    iconColor: string;
    events: Array<EventsResponse & { urgency: UrgencyMeta }>;
}

const CRITICAL_DAYS = 3;
const URGENT_DAYS = 7;
const WARNING_DAYS = 14;

export function getUrgencyLevel(daysUntil: number): UrgencyLevel {
    if (daysUntil <= CRITICAL_DAYS) return "critical";
    if (daysUntil <= URGENT_DAYS) return "urgent";
    if (daysUntil <= WARNING_DAYS) return "warning";
    return "normal";
}

export function getDayLabel(days: number): string {
    if (days < 0) return "past";
    if (days === 0) return "today";
    if (days === 1) return "tomorrow";
    return `${days} more days`;
}

export function calcDaysUntil(startTimeUnix: number): number {
    const now = new Date();
    const start = new Date(startTimeUnix * 1000);
    const diffMs =
        start.setHours(0, 0, 0, 0) - now.setHours(0, 0, 0, 0);
    return Math.floor(diffMs / (1000 * 60 * 60 * 24));
}

export function buildUrgencyMeta(event: EventsResponse): UrgencyMeta {
    const daysUntil = calcDaysUntil(event.start_time);
    return {
        level: getUrgencyLevel(daysUntil),
        daysUntil,
        label: getDayLabel(daysUntil),
    };
}

const GROUP_CONFIG: Record<
    UrgencyLevel,
    { title: string; order: number; color: string; badgeColor: string; iconColor: string }
> = {
    critical: {
        title: "Critical",
        order: 0,
        color: "text-red-600",
        badgeColor: "bg-red-100 text-red-700 border-red-200",
        iconColor: "text-red-500",
    },
    urgent: {
        title: "Urgent",
        order: 1,
        color: "text-orange-500",
        badgeColor: "bg-orange-100 text-orange-700 border-orange-200",
        iconColor: "text-orange-400",
    },
    warning: {
        title: "Warning",
        order: 2,
        color: "text-yellow-600",
        badgeColor: "bg-yellow-50 text-yellow-700 border-yellow-200",
        iconColor: "text-yellow-400",
    },
    normal: {
        title: "Normal",
        order: 3,
        color: "text-green-600",
        badgeColor: "bg-green-50 text-green-700 border-green-200",
        iconColor: "text-green-400",
    },
};

export function groupEventsByUrgency(events: EventsResponse[]): UrgencyGroup[] {
    const annotated = events.map((e) => ({ ...e, urgency: buildUrgencyMeta(e) }));

    const map = new Map<UrgencyLevel, typeof annotated>();
    for (const e of annotated) {
        const lvl = e.urgency.level;
        if (!map.has(lvl)) map.set(lvl, []);
        map.get(lvl)!.push(e);
    }

    const groups: UrgencyGroup[] = [];
    for (const [level, evts] of map.entries()) {
        const cfg = GROUP_CONFIG[level];
        groups.push({
            level,
            title: cfg.title,
            color: cfg.color,
            badgeColor: cfg.badgeColor,
            iconColor: cfg.iconColor,
            events: evts.sort((a, b) => a.urgency.daysUntil - b.urgency.daysUntil),
        });
    }

    return groups.sort(
        (a, b) => GROUP_CONFIG[a.level].order - GROUP_CONFIG[b.level].order
    );
}