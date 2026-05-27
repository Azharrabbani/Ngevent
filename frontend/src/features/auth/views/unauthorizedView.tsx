
import { useNavigate } from "react-router-dom";
import { LeftArrowIcon, ShieldIcon } from "../../../components/icon";
import Button from "../../../components/Button";

export default function UnauthorizedView() {
    const navigate = useNavigate();

    return (
        <div className="min-h-screen bg-gradient-to-br from-[#F5F7FF] to-white flex items-center justify-center px-6 py-12">
            <div className="w-full max-w-xl bg-white border border-gray-100 shadow-xl rounded-3xl p-8 sm:p-12 text-center">
                <div className="w-24 h-24 mx-auto rounded-full bg-red-100 flex items-center justify-center">
                    <ShieldIcon className="w-12 h-12 text-red-500" />
                </div>

                <div className="mt-8">
                    <h1 className="text-4xl sm:text-5xl font-bold text-gray-800">
                        403
                    </h1>

                    <h2 className="mt-3 text-2xl font-semibold text-gray-700">
                        Unauthorized Access
                    </h2>

                    <p className="mt-4 text-gray-500 leading-relaxed text-sm sm:text-base">
                        You do not have permission to access this page.
                        Please contact our administrator if you believe
                        this is a mistake.
                    </p>
                </div>

                <div className="mt-10 flex flex-col sm:flex-row items-center justify-center">
                    <Button
                        onClick={() => navigate(-1)}
                        className="w-full sm:w-auto inline-flex items-center bg-white justify-center gap-2 px-6 py-3 rounded-xl border border-gray-300 text-gray-700 font-medium hover:bg-gray-100 transition-colors"
                    >
                        <LeftArrowIcon className="w-5 h-5" />
                        Go Back
                    </Button>
                </div>
            </div>
        </div>
    );
}