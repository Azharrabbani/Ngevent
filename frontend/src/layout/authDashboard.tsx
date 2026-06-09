export default function AuthDashboard({ children }: { children: React.ReactNode }) {
    return (
        <div className="min-h-screen grid grid-cols-1 lg:grid-cols-2">
            <div className="hidden lg:block relative">
                <img
                    src="/auth_bg.png"
                    alt="Background"
                    className="absolute inset-0 w-full h-full object-cover"
                />

                <div className="absolute inset-0 bg-gradient-to-br from-blue-600/10 to-indigo-600/10" />

                <div className="relative z-10 h-full flex flex-col justify-center items-center px-10">
                    <img
                        src="/ngevent_logo.png"
                        alt="Ngevent"
                        className="w-80 h-80"
                    />
                </div>
            </div>

            <div className="flex items-center justify-center bg-white p-6">
                {children}
            </div>
        </div>
    );
}