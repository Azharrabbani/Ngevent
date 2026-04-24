import type { AttendeeResponse } from "../../profile/types/profileResponse";

interface Props {
    isOpen: boolean;
    onClose: () => void;
    profile?: AttendeeResponse;
    isLoading: boolean;
}

export default function AttendeeProfileModal({
    isOpen,
    onClose,
    profile,
    isLoading,
}: Props) {
    if (!isOpen) return null;

    return (
        <div
            className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4"
            onClick={onClose}
        >
            <div
                onClick={(e) => e.stopPropagation()}
                className="
                    w-full max-w-lg
                    max-h-[90vh] overflow-y-auto
                    bg-white border-2 border-black rounded-xl
                    p-4 sm:p-6
                    shadow-[6px_6px_0px_black]
                    relative
                "
            >
                <button
                    onClick={onClose}
                    className="absolute top-2 right-2 border-2 border-black px-2 rounded hover:bg-black hover:text-white"
                >
                    ✕
                </button>

                {isLoading ? (
                    <h1 className="text-center">Loading...</h1>
                ) : (
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
                )}
            </div>
        </div>
    );
}