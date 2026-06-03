package helper

import "time"

func BuildEventDateRange(startTime int64, month, year int) (*time.Time, *time.Time) {
	if startTime != 0 {
		loc, _ := time.LoadLocation("Asia/Jakarta")
		unix := time.Unix(startTime, 0).In(loc)

		start := time.Date(
			unix.Year(),
			unix.Month(),
			unix.Day(),
			0, 0, 0, 0,
			time.UTC,
		)

		end := start.Add(24 * time.Hour)

		return &start, &end
	}

	if month > 0 {
		if year == 0 {
			year = time.Now().Year()
		}

		start := time.Date(
			year,
			time.Month(month),
			1,
			0, 0, 0, 0,
			time.UTC,
		)

		end := start.AddDate(0, 1, 0)

		return &start, &end
	}

	if year > 0 {
		start := time.Date(
			year,
			1,
			1,
			0, 0, 0, 0,
			time.UTC,
		)

		end := start.AddDate(1, 0, 0)

		return &start, &end
	}

	return nil, nil
}
