package repository

import "ngevent/internal/dto"

type ReportRepo interface {
	GetEventReport(profileID string, filter *dto.EventReportFilter) (*dto.EventReportSummary, error)
}
