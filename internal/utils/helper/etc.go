package helper

import "time"

func ConvertDatetoUnix(date string) int64 {
	unitDate, err := time.Parse(time.RFC3339Nano, date)
	if err != nil {
		return 0
	}
	return unitDate.Unix()
}

func ConvertUnixToDate(unix int64) time.Time {
	t := time.Unix(unix, 0).UTC()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func ConvertUnixToTime(unix int64) time.Time {
	return time.Unix(unix, 0).UTC()
}
