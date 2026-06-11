export default function OrganizerSkeletonCard() {
    return (
        <div className="bg-white rounded-2xl border border-slate-200 p-4 flex items-center gap-4 animate-pulse">
            <div className="w-16 h-16 rounded-xl bg-slate-200 shrink-0" />
            <div className="flex-1 space-y-2">
                <div className="h-4 bg-slate-200 rounded w-3/4" />
                <div className="h-3 bg-slate-100 rounded w-1/2" />
                <div className="h-3 bg-slate-100 rounded w-1/4 mt-1" />
            </div>
        </div>
    );
}