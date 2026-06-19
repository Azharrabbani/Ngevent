package repository

import (
	"fmt"
	"ngevent/internal/dto"
	"ngevent/internal/model"
	"time"

	"gorm.io/gorm"
)

var wib = time.FixedZone("WIB", 7*60*60)

type ReportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) GetEventReport(profileID string, filter *dto.EventReportFilter) (*dto.EventReportSummary, error) {
	query := r.db.Model(&model.Events{}).
		Select(`
            events.*,
            organizer_profiles.name AS organizer_name
        `).
		Joins("JOIN organizer_profiles ON organizer_profiles.id = events.profile_id").
		Where("events.deleted_at IS NULL AND events.profile_id = ?", profileID)

	switch filter.Period {
	case "monthly":
		query = query.Where(
			"EXTRACT(MONTH FROM events.start_date) = ? AND EXTRACT(YEAR FROM events.start_date) = ?",
			filter.Month, filter.Year,
		)
	case "yearly":
		query = query.Where(
			"EXTRACT(YEAR FROM events.start_date) = ?",
			filter.Year,
		)
	}

	type RawRow struct {
		model.Events
		OrganizerName string `gorm:"column:organizer_name"`
	}

	var rows []RawRow
	if err := query.Scan(&rows).Error; err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	summary := &dto.EventReportSummary{
		Total: len(rows),
	}

	if filter.Period == "monthly" {
		summary.Period = fmt.Sprintf("%s %d", time.Month(filter.Month).String(), filter.Year)
	} else {
		summary.Period = fmt.Sprintf("Year %d", filter.Year)
	}

	for i, row := range rows {
		switch row.Status {
		case string(model.Active):
			summary.Active++
		case string(model.Pending):
			summary.Pending++
		case string(model.Done):
			summary.Done++
		case string(model.Rejected):
			summary.Rejected++
		}

		startTime := row.StartTime.In(wib).Format("15.04")
		endTime := row.EndTime.In(wib).Format("15.04")
		timeRange := fmt.Sprintf("%s - %s WIB", startTime, endTime)

		var rejectedReason *string
		if row.Status == string(model.Rejected) && row.RejectedReason != nil {
			rejectedReason = row.RejectedReason
		}

		summary.Rows = append(summary.Rows, dto.EventReportRow{
			No:             i + 1,
			Name:           row.Name,
			Organizer:      row.OrganizerName,
			Status:         row.Status,
			StartDate:      row.StartDate.In(wib).Format("02 Jan 2006"),
			EndDate:        row.EndDate.In(wib).Format("02 Jan 2006"),
			City:           row.City,
			TimeRange:      timeRange,
			RejectedReason: rejectedReason,
		})
	}

	return summary, nil
}
