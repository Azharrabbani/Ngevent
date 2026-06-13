export default function OwnerProfileSkeleton() {
    return (
        <div className="bg-white rounded-2xl border border-slate-200 overflow-hidden animate-pulse">
            <div className="flex flex-col md:flex-row">
                <div className="md:w-80 shrink-0 p-6 flex flex-col gap-4 md:border-r border-slate-100">
                    <div className="w-32 h-32 rounded-2xl bg-slate-200" />
                    <div className="space-y-2">
                        <div className="h-7 bg-slate-200 rounded w-3/4" />
                        <div className="h-4 bg-slate-100 rounded w-full" />
                        <div className="h-4 bg-slate-100 rounded w-5/6" />
                    </div>
                </div>

                <div className="flex-1 p-6">
                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-x-10 gap-y-6">
                        {Array.from({ length: 6 }).map((_, i) => (
                            <div key={i} className="space-y-2">
                                <div className="h-3 bg-slate-100 rounded w-1/3" />
                                <div className="h-4 bg-slate-200 rounded w-2/3" />
                            </div>
                        ))}
                    </div>
                </div>
            </div>
        </div>
    );
}