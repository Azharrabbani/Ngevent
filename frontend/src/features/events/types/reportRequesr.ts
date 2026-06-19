export interface ReportRequest {
    period: "monthly" | "yearly";
    month?: number;
    year?: number;
}