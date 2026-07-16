import { useState } from "react";
import { DownloadReportApi } from "../api/eventsApi";
import type { ReportRequest } from "../types/reportRequesr";


interface UseDownloadReportReturn {
    previewUrl: string | null;  
    isLoading: boolean;
    error: string | null;
    fetch: (params: ReportRequest) => Promise<void>;
    download: (filename: string) => void;
    reset: () => void;
}

export const useDownloadReport = (): UseDownloadReportReturn => {
    const [previewUrl, setPreviewUrl] = useState<string | null>(null);
    const [blob, setBlob] = useState<Blob | null>(null);
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const fetch = async (params: ReportRequest) => {
        if (previewUrl) {
            URL.revokeObjectURL(previewUrl);
            setPreviewUrl(null);
        }

        setIsLoading(true);
        setError(null);

        try {
            const pdfBlob = await DownloadReportApi(params);
            const url = URL.createObjectURL(pdfBlob);
            setBlob(pdfBlob);
            setPreviewUrl(url);
        } catch (err: any) {
            setError(
                err?.response?.data?.message
                ?? err?.message
                ?? "Failed to download report"
            );
        } finally {
            setIsLoading(false);
        }
    };

    const download = (filename: string) => {
        if (!blob) return;
        const url = URL.createObjectURL(blob);
        const link = document.createElement("a");
        link.href = url;
        link.download = filename;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        URL.revokeObjectURL(url);
    };

    const reset = () => {
        if (previewUrl) URL.revokeObjectURL(previewUrl);
        setPreviewUrl(null);
        setBlob(null);
        setError(null);
    };

    return { previewUrl, isLoading, error, fetch, download, reset };
};