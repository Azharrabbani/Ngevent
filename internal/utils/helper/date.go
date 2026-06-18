package helper

import "time"

func BuildEventDateRange(
	selectedDate int64,
	month int,
	year int,
) (*time.Time, *time.Time) {

	if selectedDate != 0 {

		loc, _ := time.LoadLocation("Asia/Jakarta")

		date := time.Unix(
			selectedDate,
			0,
		).In(loc)

		start := time.Date(
			date.Year(),
			date.Month(),
			date.Day(),
			0, 0, 0, 0,
			time.UTC,
		)

		end := start.Add(
			24 * time.Hour,
		)

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

		end := start.AddDate(
			0,
			1,
			0,
		)

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

		end := start.AddDate(
			1,
			0,
			0,
		)

		return &start, &end
	}

	return nil, nil
}
