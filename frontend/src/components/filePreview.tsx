interface Props {
    setPreview: () => void;
    preview: string;
};

export default function FilePreview({ setPreview, preview }: Props) {
    return (
        <div
            className="fixed inset-0 bg-black/70 flex items-center justify-center z-50"
            onClick={setPreview}
        >
            <div
                className="w-[90%] h-[90%] bg-white rounded-lg overflow-hidden"
                onClick={(e) => e.stopPropagation()}
            >
                <iframe
                    src={preview}
                    className="w-full h-full"
                />
            </div>
        </div>
    )
}