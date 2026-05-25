interface Props {
    children: React.ReactNode;
};

export default function DetailCard({ children }: Props) {
    return (
        <div className="md:max-h-[680px] md:overflow-y-auto bg-white rounded-2xl border border-gray-200 shadow-sm overflow-hidden">
            {children}
        </div>
    )
}