import { useRef } from "react";
import { cn } from "../utils/cn"

interface Props {
    className?: string
    children: React.ReactNode
    onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void
    onClickImage?: () => void;
    showEditIcon?: boolean
};

export default function UploadPhoto({
    className="",
    children,
    onChange,
    onClickImage,
    showEditIcon = false
}: Props) {
    const inputRef = useRef<HTMLInputElement>(null);

    const handleUploadClick = () => {
        inputRef.current?.click();
    };
    
    return(
        <div className="relative flex flex-col items-center">
            <input
                ref={inputRef}
                type="file"
                accept=".jpg, .jpeg, .png"
                name="photo"
                className="hidden"
                onChange={onChange} 
            />
                    
            <div
                onClick={() => {
                    if (onClickImage) {
                        onClickImage();
                    } else {
                        handleUploadClick();
                    }
                }}
                className={cn(
                    "relative flex items-center justify-center cursor-pointer transition bg-gray-200 hover:bg-gray-300",
                    className
                )}
            >
                {children}

                {showEditIcon && (
                    <div
                        onClick={(e) => {
                            e.stopPropagation();
                            handleUploadClick();
                        }}
                        className="absolute z-10 bottom-1 right-1 bg-white rounded-full p-1 shadow opacity-80 hover:opacity-100"
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" className="bi bi-pencil-square" viewBox="0 0 16 16">
                          <path d="M15.502 1.94a.5.5 0 0 1 0 .706L14.459 3.69l-2-2L13.502.646a.5.5 0 0 1 .707 0l1.293 1.293zm-1.75 2.456-2-2L4.939 9.21a.5.5 0 0 0-.121.196l-.805 2.414a.25.25 0 0 0 .316.316l2.414-.805a.5.5 0 0 0 .196-.12l6.813-6.814z"/>
                          <path fill-rule="evenodd" d="M1 13.5A1.5 1.5 0 0 0 2.5 15h11a1.5 1.5 0 0 0 1.5-1.5v-6a.5.5 0 0 0-1 0v6a.5.5 0 0 1-.5.5h-11a.5.5 0 0 1-.5-.5v-11a.5.5 0 0 1 .5-.5H9a.5.5 0 0 0 0-1H2.5A1.5 1.5 0 0 0 1 2.5z"/>
                        </svg>
                    </div>
                )}  
            </div>            
        </div>
    )
};