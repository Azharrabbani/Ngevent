import { useState } from "react";
import { useDownloadReport } from "../../hooks/useDownloadReport";
import { useGetCurrentOrganizerProfile } from "../../../profile/hooks/organizer/useGetCurrentOrganizerProfile";
import type { ReportRequest } from "../../types/reportRequesr";
import Sidebar from "../components/sidebar";
import ReportHeader from "../components/report/reportHeader";
import ErrorBanner from "../components/report/errorBanner";
import PreviewPanel from "../components/report/previewPanel";
import EmptyPreview from "../components/report/emptyPreview";
import ReportFilters from "../components/report/reportFilter";

const MONTHS = [
    "January", "February", "March", "April",
    "May", "June", "July", "August",
    "September", "October", "November", "December",
];

export default function ReportView() {
    const [period, setPeriod] = useState<"monthly" | "yearly">("monthly");
    const [month, setMonth] = useState<number>(new Date().getMonth() + 1);
    const [year, setYear] = useState<number>(new Date().getFullYear());

    const { previewUrl, isLoading, error, fetch, download, reset } = useDownloadReport();
    const organizer = useGetCurrentOrganizerProfile();

    const filename = period === "monthly"
        ? `event-report-${MONTHS[month - 1]}-${year}.pdf`
        : `event-report-${year}.pdf`;

    const periodLabel = period === "monthly"
        ? `${MONTHS[month - 1]} ${year}`
        : `Year ${year}`;

    const handlePeriodChange = (p: "monthly" | "yearly") => {
        setPeriod(p);
        reset();
    };

    const handleMonthChange = (m: number) => {
        setMonth(m);
        reset();
    };

    const handleYearChange = (y: number) => {
        setYear(y);
        reset();
    };

    const handleGenerate = () => {
        const params: ReportRequest = {
            period,
            year,
            ...(period === "monthly" ? { month } : {}),
        };
        fetch(params);
    };

    return (
        <Sidebar photoProfile={organizer?.data?.photo_profile}>
            <div className="min-h-screen bg-[#F5F6FA] p-6 md:p-10">
                <ReportHeader />

                <ReportFilters
                    period={period}
                    month={month}
                    year={year}
                    isLoading={isLoading}
                    onPeriodChange={handlePeriodChange}
                    onMonthChange={handleMonthChange}
                    onYearChange={handleYearChange}
                    onGenerate={handleGenerate}
                />

                {error && <ErrorBanner message={error} />}

                {previewUrl ? (
                    <PreviewPanel
                        previewUrl={previewUrl}
                        periodLabel={periodLabel}
                        filename={filename}
                        onDownload={() => download(filename)}
                    />
                ) : (
                    !isLoading && !error && <EmptyPreview />
                )}

            </div>
        </Sidebar>
    );
}