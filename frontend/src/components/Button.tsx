import type { ButtonHTMLAttributes } from "react";
import { cn } from "../utils/cn";

type variant = "primary" | "secondary" | "danger"

const variants = {
  primary: "bg-cyan-500 hover:bg-cyan-600",
  secondary: "bg-gray-300 hover:bg-gray-400",
  danger: "bg-red-500 hover:bg-red-600",
}

interface ButtonsProp extends ButtonHTMLAttributes<HTMLButtonElement> {
    variant?: variant
    type?:"button" | "submit";
    children: string
    className?: string
    onClick?: () => void;
}

export default function Button(
    {
        variant = "primary", 
        type = "button", 
        children,
        className = "",
        onClick
    }: ButtonsProp) {
    return(
        <button 
        type={type}
        onClick={onClick}
        className={cn(
            "text-white rounded-xl py-2 transition-all duration-300",
            variants[variant],
            className  
        )}>
            {children}
        </button>   
    )
}