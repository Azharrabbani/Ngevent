import type { AttendeeResponse } from "../../profile/types/profileResponse";

interface Props {
    profile?: AttendeeResponse;
}

export default function AttendeeProfile({ profile }: Props) {
    return (
        <>
            <div className="flex flex-col items-center gap-3 sm:gap-4">
                <img
                    src={profile?.photo_profile}
                    alt=""
                    className="w-20 h-20 sm:w-24 sm:h-24 rounded-full border-2 border-black object-cover"
                />
                <h2 className="text-lg sm:text-xl font-bold text-center">
                    {profile?.name}
                </h2>
                <p className="text-gray-600 text-sm sm:text-base">
                    @{profile?.username}
                </p>
            </div>

            <div className="mt-5 space-y-2 text-sm sm:text-base wrap-break-word">
                <p><span className="font-semibold">Email:</span> {profile?.email}</p>
                <p><span className="font-semibold">Phone:</span> {profile?.phone_number}</p>
                <p><span className="font-semibold">Country:</span> {profile?.country}</p>
                <p><span className="font-semibold">Address:</span> {profile?.address}</p>
            </div>
        </>
    )
}