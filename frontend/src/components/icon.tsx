import { BsCalendar4 } from "react-icons/bs"
import { FiClock } from "react-icons/fi"
import { LuMapPin } from "react-icons/lu";
import { AiOutlineClose } from "react-icons/ai";
import { FaPlus } from "react-icons/fa6";
import { AiOutlineCopyrightCircle } from "react-icons/ai";
import { RiCalendarEventFill } from "react-icons/ri";
import { FaRegUser } from "react-icons/fa";
import { IoHomeOutline } from "react-icons/io5";
import { ImSpinner2 } from "react-icons/im";
import { FiSearch } from "react-icons/fi";
import { GoShield } from "react-icons/go";
import { IoIosArrowRoundBack } from "react-icons/io";
import { LuUserRoundCheck } from "react-icons/lu";
import { MdOutlineDashboard } from "react-icons/md";
import { PiMaskSadFill } from "react-icons/pi";
import { ImShocked2 } from "react-icons/im";
import { FaInstagram } from "react-icons/fa";

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

export function PlusIcon({ className, size }: Props) {
    return (
        <FaPlus className={className} size={size} />
    )
}

export function CategoryIcon({ className, size }: Props) {
    return (
        <AiOutlineCopyrightCircle className={className} size={size} />
    )
}

export function EventIcon({ className, size }: Props) {
    return (
        <RiCalendarEventFill className={className} size={size} />
    )
}

export function UserIcon({ className, size }: Props) {
    return (
        <FaRegUser className={className} size={size} />
    )
}

export function HomeIcon({ className, size }: Props) {
    return (
        <IoHomeOutline className={className} size={size} />
    )
}

export function SpinnerIcon({ className, size }: Props) {
    return (
        <ImSpinner2 className={className} size={size} />
    )
}

export function SearchIcon({ className, size }: Props) {
    return (
        <FiSearch className={className} size={size} />
    )
}

export function ShieldIcon({ className, size }: Props) {
    return (
        <GoShield className={className} size={size} />
    )
}

export function LeftArrowIcon({ className, size }: Props) {
    return (
        <IoIosArrowRoundBack className={className} size={size} />
    )
}

export function UserVerifiedIcon({ className, size }: Props) {
    return (
        <LuUserRoundCheck className={className} size={size} />
    )
}

export function DashboardIcon({ className, size }: Props) {
    return (
        <MdOutlineDashboard className={className} size={size} />
    )
}

export function NoEventsIcon({ className, size }: Props) {
    return (
        <PiMaskSadFill className={className} size={size} />
    )
}

export function ShockIcon({ className, size }: Props) {
    return (
        <ImShocked2 className={className} size={size} />
    )
}

export function InstagramIcon({ className, size }: Props) {
    return (
        <FaInstagram className={className} size={size} />
    )
}