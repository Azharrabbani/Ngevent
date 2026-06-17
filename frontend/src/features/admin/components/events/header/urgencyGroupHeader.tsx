import { CircleIcon } from "../../../../../components/icon";
import type { UrgencyGroup } from "../../../utils/eventUrgency";

interface Props {
    group: UrgencyGroup;
}

export default function UrgencyGroupHeader({ group }: Props) {
    return (
        <tr className="select-none">
            <td colSpan={4} className="px-0 pt-5 pb-1">
                <div className="flex items-center gap-2 px-8">
                    <CircleIcon className={`w-2 h-2 shrink-0 ${group.iconColor}`} size={8} />

                    <span className={`text-sm font-bold ${group.color}`}>
                        {group.title}
                    </span>

                    <span className="text-xs text-gray-400 font-medium">
                        ({group.events.length})
                    </span>
                </div>
            </td>
        </tr>
    );
}