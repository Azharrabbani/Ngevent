interface ViewBannerModalProps {
    bannerUrl: string;
}

export default function ViewBannerModal({
    bannerUrl,
}: ViewBannerModalProps) {
    return (
        <div className="relative z-10 w-full max-w-4xl mx-4 animate-in fade-in zoom-in-95 duration-200">
            <div className="relative rounded-2xl overflow-hidden bg-white shadow-2xl">
                <div className="bg-black flex items-center justify-center max-h-[80vh] overflow-auto">
                    <img
                        src={bannerUrl}
                        alt="Banner Preview"
                        className="w-full object-contain"
                    />
                </div>
            </div>
        </div>
    );
}

