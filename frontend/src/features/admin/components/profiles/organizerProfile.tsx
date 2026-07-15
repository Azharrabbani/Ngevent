import { MdVerified } from "react-icons/md";
import type { OrganizerResponse, OrganizerUpdateResponse } from "../../../profile/types/profileResponse";
import OrganizerComparePreview from "./organizerComparePreview";
import OrganizerPreview from "./organizerPreview";

interface Props {
    profile?: OrganizerResponse;
    update?: OrganizerUpdateResponse;
    profileLoading?: boolean;
    updateLoading?: boolean;
    isError: boolean,
    onClose?: () => void;
};

export default function OrganizerProfile({ profile, update, profileLoading, updateLoading, isError, onClose }: Props) {
    if (!profile) {
        return null;
    }

    const hasUpdate = !!update && !isError;;
    const isApproved = profile?.status?.status === "approved";

    return (
        <>
            <div className="flex flex-col items-center gap-3 sm:gap-4">
                <img
                    src={profile.photo_profile}
                    alt=""
                    className="w-20 h-20 sm:w-24 sm:h-24 rounded-full border-2 border-black object-cover"
                />

                <div className="flex items-center gap-1">
                    <h2 className="font-bold text-lg truncate">
                        {profile.name}
                    </h2>
                    {isApproved && <MdVerified className="text-blue-600" />}
                </div>
            </div>

            {hasUpdate ? (
                <OrganizerComparePreview
                    requestUpdate={update}
                    current={profile}
                    profileLoading={profileLoading}
                    updateLoading={updateLoading}
                    onClose={onClose}
                />
            ) : (
                <OrganizerPreview
                    profile={profile}
                    loading={profileLoading}
                    onClose={onClose}
                    hasUpdate={hasUpdate}
                />
            )}
        </>
    )
}