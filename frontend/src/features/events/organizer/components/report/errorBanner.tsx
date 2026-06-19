import { WarningIcon } from "../../../../../components/icon";

interface Props {
    message: string;
}

export default function ErrorBanner({ message }: Props) {
    return (
        <div className="bg-red-50 border border-red-200 text-red-700 rounded-2xl px-5 py-4 text-sm mb-6 flex items-start gap-3">
            <WarningIcon className="text-lg mt-0.5" />
            <div>
                <p className="font-semibold mb-0.5">Failed to generate report</p>
                <p className="text-red-500">{message}</p>
            </div>
        </div>
    );
}