export default function EventViewSkeleton() {
    return (
        <div className="animate-pulse">
            <div className="grid grid-cols-1 lg:grid-cols-[1fr_360px] gap-6">
                <div className="space-y-5">
                    <div className="aspect-[16/9] rounded-2xl bg-slate-200 w-full" />

                    <div className="bg-white rounded-2xl border border-slate-200 p-6 space-y-3">
                        <div className="h-5 bg-slate-200 rounded w-40" />
                        <div className="h-3 bg-slate-100 rounded w-full" />
                        <div className="h-3 bg-slate-100 rounded w-5/6" />
                        <div className="h-3 bg-slate-100 rounded w-4/6" />
                    </div>

                    <div className="bg-white rounded-2xl border border-slate-200 p-6 space-y-4">
                        <div className="h-5 bg-slate-200 rounded w-24" />
                        <div className="h-56 rounded-xl bg-slate-100" />
                        <div className="h-3 bg-slate-100 rounded w-3/4" />
                        <div className="h-10 bg-slate-200 rounded-xl" />
                    </div>
                </div>

                <div className="space-y-4">
                    <div className="bg-white rounded-2xl border border-slate-200 p-6 space-y-5">
                        <div className="h-6 bg-slate-200 rounded w-3/4" />
                        <div className="h-3 bg-slate-100 rounded w-full" />
                        <div className="space-y-3">
                            <div className="flex gap-3">
                                <div className="w-8 h-8 rounded-lg bg-slate-200 shrink-0" />
                                <div className="space-y-1.5 flex-1">
                                    <div className="h-3.5 bg-slate-200 rounded w-3/4" />
                                    <div className="h-3 bg-slate-100 rounded w-1/2" />
                                </div>
                            </div>
                            <div className="flex gap-3">
                                <div className="w-8 h-8 rounded-lg bg-slate-200 shrink-0" />
                                <div className="space-y-1.5 flex-1">
                                    <div className="h-3.5 bg-slate-200 rounded w-2/3" />
                                    <div className="h-3 bg-slate-100 rounded w-2/5" />
                                </div>
                            </div>
                        </div>
                    </div>

                    <div className="bg-white rounded-2xl border border-slate-200 p-6 space-y-4">
                        <div className="h-3 bg-slate-100 rounded w-28" />
                        <div className="flex items-center gap-3">
                            <div className="w-12 h-12 rounded-full bg-slate-200" />
                            <div className="space-y-1.5">
                                <div className="h-3.5 bg-slate-200 rounded w-32" />
                                <div className="h-3 bg-slate-100 rounded w-20" />
                            </div>
                        </div>
                        <div className="h-10 bg-slate-200 rounded-xl" />
                    </div>
                </div>
            </div>
        </div>
    );
}