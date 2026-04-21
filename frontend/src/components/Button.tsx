import type { ButtonHTMLAttributes } from "react";
import { cn } from "../utils/cn";

type variant = "primary" | "secondary" | "danger"

const variants = {
  primary: "bg-[#0040A1] hover:bg-blue-700",
  secondary: "bg-gray-300 hover:bg-gray-400",
  danger: "bg-red-500 hover:bg-red-600",
}

interface ButtonsProp extends ButtonHTMLAttributes<HTMLButtonElement> {
    variant?: variant
    children: string
    className?: string
}

export default function Button(
    {
        variant = "primary", 
        children,
        className = "",
        onClick,
        ...props
    }: ButtonsProp) {
    return(
        <button
        type="button" 
        {...props}
        onClick={onClick}
        className={cn(
          "text-white rounded-xl py-2 transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed",
          variants[variant],
          className  
        )}>
            {children}
        </button>   
    )
}