import { useEffect, useState } from "react"

interface Props {
    preview: string
    setPreview: () => void
}

export default function FilePreview({ preview, setPreview }: Props) {
    const [blobUrl, setBlobUrl] = useState<string | null>(null)
    const [loading, setLoading] = useState(true)

    useEffect(() => {        
        if (preview.startsWith("blob:")) {
            setBlobUrl(preview)
            setLoading(false)
            return
        }
    
        setLoading(true)
        fetch(preview, { credentials: "include" })
            .then(res => res.blob())
            .then(blob => {
                setBlobUrl(URL.createObjectURL(blob))
                setLoading(false)
            })
            .catch(() => setLoading(false))

        return () => {
            if (blobUrl && !blobUrl.startsWith("blob:")) {
                URL.revokeObjectURL(blobUrl)
            }
        }
    }, [preview])

    return (
        <div
            className="fixed inset-0 bg-black/70 flex items-center justify-center z-50"
            onClick={setPreview}
        >
            <div
                className="w-[90%] h-[90%] bg-white rounded-lg overflow-hidden"
                onClick={(e) => e.stopPropagation()}
            >
                {loading ? (
                    <div className="flex items-center justify-center h-full text-gray-500">
                        Loading...
                    </div>
                ) : blobUrl ? (
                    <iframe src={blobUrl} className="w-full h-full" />
                ) : (
                    <div className="flex items-center justify-center h-full text-red-500">
                        Failed to load file
                    </div>
                )}
            </div>
        </div>
    )
}