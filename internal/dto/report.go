package dto

type EventReportFilter struct {
	Period string `json:"period" query:"period"`
	Month  int    `json:"month"  query:"month"`
	Year   int    `json:"year"   query:"year"`
}

type EventReportRow struct {
	No             int
	Name           string
	Organizer      string
	Status         string
	StartDate      string
	EndDate        string
	City           string
	TimeRange      string 
	RejectedReason *string
}

type EventReportSummary struct {
	Period   string
	Total    int
	Active   int
	Pending  int
	Done     int
	Rejected int
	Rows     []EventReportRow
}