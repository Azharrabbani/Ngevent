package helper

import "time"

func ConvertDatetoUnix(date string) int64 {
	unitDate, err := time.Parse(time.RFC3339Nano, date)
	if err != nil {
		return 0
	}

	return unitDate.Unix()
}

func ConvertUnixtoDate(unix int64) time.Time {
	return time.Unix(unix, 0)
}
