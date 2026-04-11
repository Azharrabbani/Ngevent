import { useState } from "react";
import AuthContainer from "../components/container";
import OtpInput from "../components/otpInput";
import Button from "../../../components/Button";

export default function VerifiedEmailView() {
    const [value, setValue] = useState("");

    return (
        <AuthContainer className="p-20 flex flex-col items-center gap-7">        
            <div className="w-fit p-6 bg-blue-100 rounded-2xl">
                <svg xmlns="http://www.w3.org/2000/svg" width="38" height="38" fill="currentColor" className="bi bi-shield-lock-fill text-blue-600" viewBox="0 0 16 16">
                <path fill-rule="evenodd" d="M8 0c-.69 0-1.843.265-2.928.56-1.11.3-2.229.655-2.887.87a1.54 1.54 0 0 0-1.044 1.262c-.596 4.477.787 7.795 2.465 9.99a11.8 11.8 0 0 0 2.517 2.453c.386.273.744.482 1.048.625.28.132.581.24.829.24s.548-.108.829-.24a7 7 0 0 0 1.048-.625 11.8 11.8 0 0 0 2.517-2.453c1.678-2.195 3.061-5.513 2.465-9.99a1.54 1.54 0 0 0-1.044-1.263 63 63 0 0 0-2.887-.87C9.843.266 8.69 0 8 0m0 5a1.5 1.5 0 0 1 .5 2.915l.385 1.99a.5.5 0 0 1-.491.595h-.788a.5.5 0 0 1-.49-.595l.384-1.99A1.5 1.5 0 0 1 8 5"/>
                </svg>
            </div>
        
            <div className="text-center flex flex-col gap-3">
                <h1 className="font-bold text-4xl">
                    Verification Code
                </h1>
                <p className="max-w-sm text-center">
                    Please enter the OTP code sent to your email to verify your account and 
                    continue using our event platform securely.
                </p>
            </div>
            
            <OtpInput
                digitLength={6}
                value={value}
                onChange={e => setValue(e.target.value)}
                
            />

            <Button
                
                className="bg-blue-600 hover:bg-blue-500 w-full text-xl"
            >
                Verify
            </Button>

            <div className="flex justify-between w-full">
                <Button
                    className="bg-gray-300 hover:bg-blue-300 hover:text-white font-bold text-sm px-5 text-gray-800"
                >
                    Resend
                </Button>

                <p 
                    className="flex items-center gap-2 bg-blue-300 text-center font-bold text-white rounded-xl p-2 px-5 text-sm"
                >
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" className="bi bi-stopwatch-fill" viewBox="0 0 16 16">
                      <path d="M6.5 0a.5.5 0 0 0 0 1H7v1.07A7.001 7.001 0 0 0 8 16a7 7 0 0 0 5.29-11.584l.013-.012.354-.354.353.354a.5.5 0 1 0 .707-.707l-1.414-1.415a.5.5 0 1 0-.707.707l.354.354-.354.354-.012.012A6.97 6.97 0 0 0 9 2.071V1h.5a.5.5 0 0 0 0-1zm2 5.6V9a.5.5 0 0 1-.5.5H4.5a.5.5 0 0 1 0-1h3V5.6a.5.5 0 1 1 1 0"/>
                    </svg>
                   30s 
                </p>
            </div>
            

        
            

        </AuthContainer>
    )
}