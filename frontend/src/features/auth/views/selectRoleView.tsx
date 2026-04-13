import { useState } from "react";
import Button from "../../../components/Button";
import AuthContainer from "../components/container";
import { useSelectRole } from "../hooks/useSelectRole";

export default function SelectRoleView() {
    const [role, setRole] = useState("");

    const {selectRole, loading, message, error} = useSelectRole();

    const handleSelectRole = async(e: React.MouseEvent) => {
        e.preventDefault();

        if (!role) {
            alert("Please select a role first");
            return;
        }

        await selectRole({role: role});
    }

    return(
        <AuthContainer className="flex flex-col md:flex md:flex-col gap-10 p-8">
            <div className="flex flex-col text-center w-fit gap-2">
                <h1 className="text-3xl font-bold">Please select your role</h1>
                <p className="max-w-md mx-auto">Select your role to unlock features and events that match your needs.</p>
            </div>

            <div className=" md:flex md:gap-10">
                {/* Attendee box */}
                <div
                onClick={() => setRole("user")}
                className={`
                    flex flex-col gap-3 bg-gray-100 rounded-xl max-w-96 p-10 
                    hover:ring-3 hover:ring-blue-400 hover:cursor-pointer 
                    transition-all duration-200
                    ${role === "user" ? "ring-3 ring-blue-400" : ""}
                `}>
                    <div className="w-fit bg-blue-200 rounded-full p-3">
                        <svg xmlns="http://www.w3.org/2000/svg" width="30" height="25" fill="currentColor" className="bi bi-people text-[#1A146B]" viewBox="0 0 16 16">
                          <path d="M15 14s1 0 1-1-1-4-5-4-5 3-5 4 1 1 1 1zm-7.978-1L7 12.996c.001-.264.167-1.03.76-1.72C8.312 10.629 9.282 10 11 10c1.717 0 2.687.63 3.24 1.276.593.69.758 1.457.76 1.72l-.008.002-.014.002zM11 7a2 2 0 1 0 0-4 2 2 0 0 0 0 4m3-2a3 3 0 1 1-6 0 3 3 0 0 1 6 0M6.936 9.28a6 6 0 0 0-1.23-.247A7 7 0 0 0 5 9c-4 0-5 3-5 4q0 1 1 1h4.216A2.24 2.24 0 0 1 5 13c0-1.01.377-2.042 1.09-2.904.243-.294.526-.569.846-.816M4.92 10A5.5 5.5 0 0 0 4 13H1c0-.26.164-1.03.76-1.724.545-.636 1.492-1.256 3.16-1.275ZM1.5 5.5a3 3 0 1 1 6 0 3 3 0 0 1-6 0m3-2a2 2 0 1 0 0 4 2 2 0 0 0 0-4"/>
                        </svg>
                    </div>

                    <div>
                        <h1 className="text-xl font-semibold mb-2">Attendee</h1>
                        <p className="max-w-sm mx-auto">
                            Discover events and expand your network
                            through curated experiences tailored to your
                            interests.
                        </p>
                    </div>
                </div>

                {/* Organizer box */}
                <div 
                onClick={() => setRole("event organizer")}
                className={`
                    flex flex-col gap-3 bg-gray-100 rounded-xl max-w-96 p-7 
                    hover:ring-3 hover:ring-[#312E81] hover:cursor-pointer 
                    transition-all duration-200 mt-7
                    ${role === "event organizer" ? "ring-3 ring-[#312E81]" : ""}
                `}>
                    <div className="w-fit rounded-full p-3 bg-[#312E81]">
                        <svg xmlns="http://www.w3.org/2000/svg" width="30" height="25" fill="currentColor" className="bi bi-calendar text-[#9C9AF4]" viewBox="0 0 16 16">
                          <path d="M3.5 0a.5.5 0 0 1 .5.5V1h8V.5a.5.5 0 0 1 1 0V1h1a2 2 0 0 1 2 2v11a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2V3a2 2 0 0 1 2-2h1V.5a.5.5 0 0 1 .5-.5M1 4v10a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1V4z"/>
                        </svg>
                    </div>

                    <div>
                        <h1 className="text-xl font-semibold mb-2">Event Organizer</h1>
                        <p className="max-w-sm mx-auto">
                            Create, manage, and grow your events with our
                            powerful suite of management tools and
                            analytics.
                        </p>
                    </div>
                </div>
            </div>
            
            {error && (
                <p className="text-red-500 text-md text-center">
                    {error}
                </p>
            )}

            {message && (
                <p className="text-green-600 text-md text-center">
                    {message}
                </p>
             )}

            <Button
            onClick={handleSelectRole}
            disabled={loading || role === ""}
           className={`
                p-2 w-64 rounded-md
                ${!role 
                    ? "bg-gray-300 cursor-not-allowed" 
                    : role === "event organizer"
                    ? "bg-[#312E81] hover:bg-purple-900 text-white"
                    : "bg-cyan-500 hover:bg-cyan-600 text-white"}
            `}>
                {loading ? "Loading..." : "Continue"}
            </Button>
        </AuthContainer>
    )
}