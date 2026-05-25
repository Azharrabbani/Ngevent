import { BsCalendar4 } from "react-icons/bs"
import { FiClock } from "react-icons/fi"
import { LuMapPin } from "react-icons/lu";
import { AiOutlineClose } from "react-icons/ai";

interface Props {
    className?: string
    size?: number
}

export function CalenderIcon({ className, size }: Props) {
    return (
        <BsCalendar4 className={className} size={size} />
    )
}

export function ClockIcon({ className, size }: Props) {
    return (
        <FiClock className={className} size={size} />
    )
}

export function PinIcon({ className, size }: Props) {
    return (
        <LuMapPin className={className} size={size} />
    )
}

export function CrossIcon({ className, size }: Props) {
    return (
        <AiOutlineClose className={className} size={size} />
    )
}