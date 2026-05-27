interface Props {
    children: React.ReactElement;
};

export default function EventFormContainer( { children }: Props ) {
    return (
        <div className="min-h-screen flex items-center justify-center bg-[#F4F7FB]">
            {children}
        </div>
    )
}