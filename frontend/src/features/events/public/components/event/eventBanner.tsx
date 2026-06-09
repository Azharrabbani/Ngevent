import type { EventCategory } from "../../../types/publicEventResponse";

interface Props {
    banner: string;
    name: string;
    categories: EventCategory[];
}

export default function EventBanner({ banner, name, categories }: Props) {
    return (
        <div className="relative rounded-2xl overflow-hidden bg-slate-100 aspect-[16/9] w-full">
            <img
                src={banner}
                alt={name}
                className="w-full h-full object-cover"
            />
            <div className="absolute inset-0 bg-gradient-to-t from-black/20 to-transparent" />

            {categories.length > 0 && (
                <div className="absolute top-3 left-3 flex flex-wrap gap-2">
                    {categories.map((cat) => (
                        <span
                            key={cat.id}
                            className="inline-flex items-center gap-1.5 bg-white/90 backdrop-blur-sm text-slate-700 text-xs font-medium px-3 py-1.5 rounded-full shadow-sm"
                        >
                            <span className="text-[10px]">✦</span>
                            {cat.name}
                        </span>
                    ))}
                </div>
            )}
        </div>
    );
}