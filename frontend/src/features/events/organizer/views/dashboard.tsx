import Sidebar from "../components/sidebar";
import Header from "../components/header";
import EventCard from "../components/eventCard";
import Pagination from "../../../../components/pagination";
import { useState } from "react";

export default function Dashboard() {
    const [currentPage, setCurrentPage] = useState(1);

    const totalPage = 20;

    return (
        <Sidebar>
            <>
                <Header/>
                
                <div className="p-12">
                    <div className="space-y-2 text-center md:text-start">
                        <h2 className="text-lg font-semibold">Active Events</h2>
                        <p className="text-[#424654]">Showing xxx events</p>
                    </div>

                    <div className="mt-8 grid grid-cols-1 sm:grid-cols-2 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
                        <EventCard
                            title="Event title"
                            date="Oct 15 - 18, 2024"
                            location="San Francisco"
                            status="active"
                            image="https://img.freepik.com/free-photo/abstract-flowing-neon-wave-background_53876-101942.jpg"
                            revenue={124000}
                            tickets={[
                                { name: "VIP Access", price: 120000, sold: 170, total: 200 },
                                { name: "Regular", price: 50000, sold: 80, total: 150 }
                            ]}
                        />                    
                        <EventCard
                            title="Event title"
                            date="Oct 15 - 18, 2024"
                            location="San Francisco"
                            status="active"
                            image="https://img.freepik.com/free-photo/abstract-flowing-neon-wave-background_53876-101942.jpg"
                            revenue={124000}
                            tickets={[
                                { name: "VIP Access", price: 120000, sold: 170, total: 200 },
                                { name: "Regular", price: 50000, sold: 80, total: 150 }
                            ]}
                        />                    
                        <EventCard
                            title="Event title"
                            date="Oct 15 - 18, 2024"
                            location="San Francisco"
                            status="active"
                            image="https://img.freepik.com/free-photo/abstract-flowing-neon-wave-background_53876-101942.jpg"
                            revenue={124000}
                            tickets={[
                                { name: "VIP Access", price: 120000, sold: 170, total: 200 },
                                { name: "Regular", price: 50000, sold: 80, total: 150 }
                            ]}
                        />                    
                        <EventCard
                            title="Event title"
                            date="Oct 15 - 18, 2024"
                            location="San Francisco"
                            status="active"
                            image="https://img.freepik.com/free-photo/abstract-flowing-neon-wave-background_53876-101942.jpg"
                            revenue={124000}
                            tickets={[
                                { name: "VIP Access", price: 120000, sold: 170, total: 200 },
                                { name: "Regular", price: 50000, sold: 80, total: 150 }
                            ]}
                        />                    
                    </div>
                </div>

                <Pagination
                    currentPage={currentPage}
                    totalPage={totalPage}
                    onPrev={() => setCurrentPage((prev) => Math.max(prev - 1, 1))}
                    onNext={() => setCurrentPage((prev) => Math.min(prev + 1, totalPage))}
                    onCurrent={(page) => setCurrentPage(page)}
                />

            </>
        </Sidebar>
    )
}