import type { UrgencyMeta } from "../../../utils/eventUrgency";


interface Props {
    urgency: UrgencyMeta;
    badgeColor: string;
}

export default function UrgencyBadge({ urgency, badgeColor }: Props) {
    return (
        <span
            className={`
                inline-flex items-center gap-1.5 px-2.5 py-1
                rounded-full text-xs font-medium border
                whitespace-nowrap
                ${badgeColor}
            `}
        >
            {urgency.label}
        </span>
    );
}