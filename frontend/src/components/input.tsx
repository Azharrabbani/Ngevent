import { forwardRef, type InputHTMLAttributes, type ReactNode } from "react";
import { cn } from "../utils/cn";

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  labelStyle?: string;
  error?: string;
  leftIcon?: ReactNode;
  rightIcon?: ReactNode;
}

const Input = forwardRef<HTMLInputElement, InputProps>(({
  label,
  type = "text",
  name,
  placeholder,
  className = "",
  labelStyle = "",
  error,
  leftIcon,
  rightIcon,
  ...rest
}, ref) => {
  return (
    <div className="w-full">
      {label && (
        <label className={cn("block text-sm mb-1 font-medium", labelStyle)}>
          {label}
        </label>
      )}

      <div className="relative">
        {leftIcon && (
          <div className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400">
            {leftIcon}
          </div>
        )}

        <input
          ref={ref}
          type={type}
          name={name}
          placeholder={placeholder}
          {...rest}
          className={cn(
            "w-full p-2 rounded-xl bg-gray-200 outline-none",
            leftIcon && "pl-10",
            rightIcon && "pr-10",
            error && "border border-red-500",
            className,
          )}
        />

        {rightIcon && (
          <div className="absolute right-3 top-1/2 -translate-y-1/2 cursor-pointer text-gray-400">
            {rightIcon}
          </div>
        )}
      </div>

      {error && (
        <p className="text-red-500 text-xs mt-1">{error}</p>
      )}
    </div>
  );
});

Input.displayName = "Input";

export default Input;