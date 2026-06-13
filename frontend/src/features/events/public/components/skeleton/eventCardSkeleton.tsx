export default function EventCardSkeleton() {
    return (
        <div className="bg-white rounded-2xl border border-slate-200 overflow-hidden animate-pulse">
            <div className="h-44 bg-slate-200" />
            <div className="p-4 space-y-3">
                <div className="h-5 bg-slate-200 rounded w-3/4" />
                <div className="space-y-1.5 mt-3">
                    <div className="h-3 bg-slate-100 rounded w-2/5" />
                    <div className="h-3 bg-slate-100 rounded w-1/2" />
                </div>
                <div className="pt-3 border-t border-slate-100 flex items-center gap-2">
                    <div className="w-6 h-6 rounded-full bg-slate-200 shrink-0" />
                    <div className="h-3 bg-slate-100 rounded w-1/3" />
                </div>
            </div>
        </div>
    );
}