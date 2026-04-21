import { useRef } from "react";

interface Props {
    uniqueId: string
    children: string
    file: string | undefined
    onChange?: (file: File) => void
    onClickFile?: () => void;
    showEditIcon?: boolean
}

export default function UploadFile({
    uniqueId,
    file,
    children,
    onChange,
    onClickFile,
    showEditIcon = true,
}: Props) {
    const inputRef = useRef<HTMLInputElement>(null);

    return (
        <div className="flex flex-col items-center gap-2 w-full">
            <p className="text-xs tracking-widest text-center">
                {children}
            </p>

            <input
                ref={inputRef}
                id={uniqueId}
                type="file"
                accept="application/pdf"
                className="hidden"
                onChange={(e) => {
                    const file = e.target.files?.[0];
                    if (!file) return;
                    onChange?.(file);
                }}
            />

            <div className="relative w-full">
                <div
                    onClick={onClickFile}
                    className="w-full h-28 sm:h-32
                               rounded-xl border-2 border-dashed border-gray-400
                               bg-gray-100
                               flex flex-col gap-2 items-center justify-center
                               cursor-pointer hover:bg-gray-200 transition"
                >
                    <svg xmlns="http://www.w3.org/2000/svg" width="25" height="25" fill="red" className="bi bi-filetype-pdf" viewBox="0 0 16 16">
                        <path fill-rule="evenodd" d="M14 4.5V14a2 2 0 0 1-2 2h-1v-1h1a1 1 0 0 0 1-1V4.5h-2A1.5 1.5 0 0 1 9.5 3V1H4a1 1 0 0 0-1 1v9H2V2a2 2 0 0 1 2-2h5.5zM1.6 11.85H0v3.999h.791v-1.342h.803q.43 0 .732-.173.305-.175.463-.474a1.4 1.4 0 0 0 .161-.677q0-.375-.158-.677a1.2 1.2 0 0 0-.46-.477q-.3-.18-.732-.179m.545 1.333a.8.8 0 0 1-.085.38.57.57 0 0 1-.238.241.8.8 0 0 1-.375.082H.788V12.48h.66q.327 0 .512.181.185.183.185.522m1.217-1.333v3.999h1.46q.602 0 .998-.237a1.45 1.45 0 0 0 .595-.689q.196-.45.196-1.084 0-.63-.196-1.075a1.43 1.43 0 0 0-.589-.68q-.396-.234-1.005-.234zm.791.645h.563q.371 0 .609.152a.9.9 0 0 1 .354.454q.118.302.118.753a2.3 2.3 0 0 1-.068.592 1.1 1.1 0 0 1-.196.422.8.8 0 0 1-.334.252 1.3 1.3 0 0 1-.483.082h-.563zm3.743 1.763v1.591h-.79V11.85h2.548v.653H7.896v1.117h1.606v.638z"/>
                    </svg>

                    <p className="text-xs text-gray-600">
                        {file ? "View File" : "Upload File"}
                    </p>
                </div>

                {showEditIcon && (
                    <label
                        htmlFor={uniqueId}
                        onClick={(e) => e.stopPropagation()}
                        className="absolute bottom-2 right-2 bg-white rounded-full p-2 shadow cursor-pointer"
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" viewBox="0 0 16 16">
                            <path d="M15.502 1.94a.5.5 0 0 1 0 .706L14.459 3.69l-2-2L13.502.646a.5.5 0 0 1 .707 0l1.293 1.293z"/>
                            <path d="M4.939 9.21l6.813-6.814 2 2-6.813 6.814a.5.5 0 0 1-.196.12l-2.414.805a.25.25 0 0 1-.316-.316l.805-2.414a.5.5 0 0 1 .121-.196z"/>
                        </svg>
                    </label>
                )}
            </div>
        </div>
    );
}