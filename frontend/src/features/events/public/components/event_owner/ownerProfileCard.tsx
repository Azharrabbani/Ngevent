import { FiMail, FiMapPin, FiPhone, FiCalendar, FiCheckCircle } from "react-icons/fi";
import type { OrganizerResponse } from "../../../../profile/types/profileResponse";
import { InstagramIcon } from "../../../../../components/icon";

interface Props {
    organizer: OrganizerResponse;
}

interface InfoRowProps {
    icon: React.ReactNode;
    label: string;
    value: string;
    highlight?: boolean;
}

function InfoRow({ icon, label, value, highlight }: InfoRowProps) {
    return (
        <div className="flex flex-col gap-0.5">
            <div className="flex items-center gap-2 text-slate-400">
                <span className="w-4 h-4 shrink-0">{icon}</span>
                <span className="text-xs">{label}</span>
            </div>
            <p className={`text-sm font-medium pl-6 ${highlight ? "text-blue-600" : "text-slate-800"}`}>
                {value || "—"}
            </p>
        </div>
    );
}

interface InstagramRowProps {
    instagram: string | undefined;
}

function InstagramRow({ instagram }: InstagramRowProps) {
    const href = instagram
        ? instagram.startsWith("http")
            ? instagram
            : `https://www.instagram.com/${instagram.replace(/^@/, "")}`
        : null;

    return (
        <div className="flex flex-col gap-1.5">
            {href && (
                <a
                    href={href}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="w-fit"
                >
                    <InstagramIcon className="text-3xl rounded-xl bg-gradient-to-r from-pink-500 to-purple-500 text-white hover:opacity-90 transition-opacity" />
                </a>
            )}
        </div>
    );
}

export default function OrganizerProfileCard({ organizer }: Props) {
    const joinedYear = organizer.company_detail
        ? new Date().getFullYear()
        : 2024;

    return (
        <div className="bg-white rounded-2xl border border-slate-200 overflow-hidden">
            <div className="flex flex-col md:flex-row">
                <div className="md:w-72 shrink-0 p-6 flex flex-col gap-4 md:border-r border-slate-100">
                    <div className="w-32 h-32 rounded-2xl overflow-hidden border-2 border-slate-100 shrink-0">
                        {organizer.photo_profile ? (
                            <img
                                src={organizer.photo_profile}
                                alt={organizer.name}
                                className="w-full h-full object-cover"
                            />
                        ) : (
                            <div className="w-full h-full bg-gradient-to-br from-blue-100 to-indigo-200 flex items-center justify-center text-blue-600 font-bold text-4xl">
                                {organizer.name?.charAt(0)?.toUpperCase() ?? "?"}
                            </div>
                        )}
                    </div>

                    <div className="flex items-center gap-2 flex-wrap">
                        <h1 className="text-2xl font-bold text-slate-900">
                            {organizer.name}
                        </h1>
                        {organizer.status.status === "approved" && (
                            <FiCheckCircle className="w-5 h-5 text-blue-500 shrink-0" />
                        )}
                    </div>

                    <InstagramRow instagram={organizer.social_media?.instagram} />
                </div>

                <div className="flex-1 p-6 flex flex-col gap-6">
                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-x-10 gap-y-6">
                        <InfoRow
                            icon={<FiMail className="w-4 h-4" />}
                            label="Email"
                            value={organizer.social_media?.email || organizer.email}
                            highlight
                        />
                        <InfoRow
                            icon={<FiMail className="w-4 h-4" />}
                            label="Secondary Email"
                            value={organizer.social_media?.email}
                            highlight
                        />
                        <InfoRow
                            icon={<FiPhone className="w-4 h-4" />}
                            label="Phone"
                            value={organizer.phone_number}
                        />
                        <InfoRow
                            icon={<FiMapPin className="w-4 h-4" />}
                            label="Country"
                            value={organizer.country}
                            highlight
                        />
                        <InfoRow
                            icon={<FiCalendar className="w-4 h-4" />}
                            label="Joined since"
                            value={String(joinedYear)}
                        />
                    </div>
                    {organizer.company_detail?.description && (
                        <div className="flex flex-col gap-1.5 border-t border-slate-100 pt-5">
                            <span className="text-xs text-slate-400">About</span>
                            <div
                                className="
                                        min-h-48
                                        max-h-72
                                        overflow-y-auto
                                        text-sm text-slate-600 leading-relaxed
                                        pr-2
                                    "
                                dangerouslySetInnerHTML={{
                                    __html: organizer.company_detail.description,
                                }}
                            />
                        </div>
                    )}
                </div>

            </div>
        </div>
    );
}