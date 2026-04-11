import React, { useState, type ComponentProps, useRef, useId, useEffect } from "react"
import { cn } from "../../../utils/cn";

type PropsType = Omit<ComponentProps<"input">, "value"> & {
    digitLength?: 4 | 6;
    value?: string;
};

export default function OtpInput({
    digitLength = 6,
    className,
    disabled,
    onChange,
    value,
    ...props
}: PropsType) {
    const [otp, setOtp] = useState(Array(digitLength).fill(""));
    const inputRefs = useRef<(HTMLInputElement | null) []>([]);

    const id = useId()

    useEffect(() => {
        if (value !== undefined) {
            const digits = value.split("").slice(0, digitLength);
            const paddedDiggits = [
                ...digits,
                ...Array(digitLength - digits.length).fill("")
            ];
            setOtp(paddedDiggits);
        }
    }, [value, digitLength]);

    useEffect(() => {
        if (onChange) {
        const otpValue = otp.join("");
        const event = {
            target: { value: otpValue }
        } as React.ChangeEvent<HTMLInputElement>;
        onChange(event);
        }
    }, [otp, onChange]);

    const handleKeyDown = (e : React.KeyboardEvent<HTMLInputElement>) => {
        if (
            !/^[0-9]{1}$/.test(e.key) &&
            e.key !== "Backspace" &&
            e.key !== "Delete" &&
            e.key !== "Tab" &&
            !e.metaKey
        ) {
            e.preventDefault()
        }

        if (e.key === "Delete" || e.key === "Backspace") {
            const index = inputRefs.current.indexOf(e.target as HTMLInputElement);

            if (otp[index]) {
                setOtp(prevOtp => [
                    ...prevOtp.slice(0, index),
                    "",
                    ...prevOtp.slice(index + 1)
                ]);
            } else if (index > 0) {
                setOtp(prevOtp => [
                    ...prevOtp.slice(0, index - 1),
                    "",
                    ...prevOtp.slice(index)
                ]);
                inputRefs.current[index - 1]?.focus();
            }
        }
    };

    const handleInput = (e: React.ChangeEvent<HTMLInputElement>) => {
        const {target} = e;
        const index = inputRefs.current.indexOf(target as HTMLInputElement);

        if (target.value) {
            setOtp(prevOtp => [
                ...prevOtp.slice(0, index),
                target.value,
                ...prevOtp.slice(index + 1)
            ]);

            if (index < otp.length - 1) {
                inputRefs.current[index + 1]?.focus();
            }
        }
    };

    const handleFocus = (e: React.FocusEvent<HTMLInputElement, Element>) => {
        e.target.select();
    };

    const handlePaste = (e: React.ClipboardEvent<HTMLInputElement>) => {
        e.preventDefault();
        const text = e.clipboardData.getData("text");
        if (!new RegExp(`^[0-9]{${otp.length}}$`).test(text)) {
            return;
        }

        const digits = text.split("");
        setOtp(digits);
    };

    return (
        <div>
            <div className="flex items-center gap-2 not-focus-within:text-input-placeholder-text focus-within:text-text-50">
                {otp.map((digit, index) => (
                    <input 
                        type="text"
                        key={index}
                        id={index === 0 ? id : undefined}
                        maxLength={1}
                        value={digit}
                        onChange={handleInput}
                        onKeyDown={handleKeyDown}
                        onFocus={handleFocus}
                        onPaste={handlePaste}
                        ref={el => {
                            inputRefs.current[index] = el;
                        }}
                        className={cn(
                            "bg-gray-300 focus:border-input-primary-focus-border focus:ring-input-primary-focus-border/20 flex size-15 items-center justify-center rounded-lg p-2 text-center text-base font-normal shadow-xs focus:ring-1 focus:outline-none disabled:border-base-100 disabled:bg-background-soft-50 disabled:text-input-disabled-text",
                            digit && "text-text-50",
                            className
                        )}
                        style={{
                            order: digitLength === 6 && index > 2 ? index + 1: index
                        }}
                        disabled={disabled}
                        {...props}
                    />
                ))}

                
            </div>
        </div>
    );
}