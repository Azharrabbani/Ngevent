import { forwardRef, type InputHTMLAttributes, type ReactNode } from "react";
import { cn } from "../utils/cn";

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  labelStyle?: string;
  error?: string;
  leftIcon?: ReactNode;
  rightIcon?: ReactNode;
  onlyNumber?: boolean;
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
  onlyNumber = false,
  onChange,
  onKeyDown,
  ...rest
}, ref) => {
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (onlyNumber) {
      e.target.value = e.target.value.replace(/[^0-9]/g, "");
    }
    onChange?.(e);
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (onlyNumber) {
      const allowedKeys = [
        "Backspace", "Delete", "Tab", "Escape", "Enter",
        "ArrowLeft", "ArrowRight", "ArrowUp", "ArrowDown",
        "Home", "End",
      ];
      const isCtrlCombo = e.ctrlKey || e.metaKey;
      const isDigit = /^[0-9]$/.test(e.key);

      if (!isDigit && !allowedKeys.includes(e.key) && !isCtrlCombo) {
        e.preventDefault();
      }
    }
    onKeyDown?.(e);
  };

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
          type={onlyNumber ? "text" : type}
          inputMode={onlyNumber ? "numeric" : rest.inputMode}
          pattern={onlyNumber ? "[0-9]*" : rest.pattern}
          name={name}
          placeholder={placeholder}
          onChange={handleChange}
          onKeyDown={handleKeyDown}
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