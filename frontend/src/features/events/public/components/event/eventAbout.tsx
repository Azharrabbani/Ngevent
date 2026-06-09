interface Props {
    description: string;
}

export default function EventAbout({ description }: Props) {
    return (
        <div className="bg-white rounded-2xl border border-slate-200 p-6">
            <h2 className="text-lg font-bold text-slate-900 mb-4">About this Event</h2>
            <div
                className="prose prose-sm text-sm text-gray-600 leading-relaxed space-y-3"
                dangerouslySetInnerHTML={{ __html: description }}
            />
        </div>
    );
}