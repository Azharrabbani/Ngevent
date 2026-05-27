import { useState } from "react";
import ViewBannerModal from "./bannerModal";
import Modal from "../../../../../components/modal";


interface BannerPreviewProps {
    imageUrl: string;
    title: string;
}

export default function BannerPreview({
    imageUrl,
    title,
}: BannerPreviewProps) {
    const [open, setOpen] = useState(false);

    return (
        <>
            <img
                src={imageUrl}
                alt={title}
                onClick={() => setOpen(true)}
                className="w-full h-56 object-cover rounded-lg cursor-pointer hover:opacity-90 transition"
            />

            <Modal
                isOpen={open}
                onClose={() => setOpen(false)}
            >
                <ViewBannerModal
                    bannerUrl={imageUrl}
                />
            </Modal>
        </>
    );
}