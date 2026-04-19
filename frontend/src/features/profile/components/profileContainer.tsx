import type React from "react"

interface Props{
    children: React.ReactNode
    role: string
};

export default function ProfileContainer({children, role}: Props) {
    return(
        <section className="bg-gray-100 min-h-screen 
                            grid grid-cols-1 
                            md:grid md:grid-cols-[200px_1fr_200px] ">
            <div className={`hidden md:block
                            ${role === "user" ? 
                            "bg-[url('https://png.pngtree.com/background/20230614/original/pngtree-abstract-blue-wave-background-vector-picture-image_3513788.jpg')]"
                            : 
                            "bg-[url('https://img.freepik.com/premium-vector/gradient-abstract-purple-background_23-2151877363.jpg?semt=ais_hybrid&w=740&q=80')]"} 
                             
                            bg-cover bg-center`}/>

            <div className="flex items-center justify-center px-4 py-8">
                <div className="w-full max-w-2xl">
                    {children}
                </div>
            </div>

            <div className={`hidden md:block
                            ${role === "event organizer" ? 
                            "bg-[url('https://img.freepik.com/premium-vector/gradient-abstract-purple-background_23-2151877363.jpg?semt=ais_hybrid&w=740&q=80')]"
                            : 
                            "bg-[url('https://png.pngtree.com/background/20230614/original/pngtree-abstract-blue-wave-background-vector-picture-image_3513788.jpg')]"} 
                             
                            bg-cover bg-center`}/>
        </section>
    )
}