package helper

func GetEventStatusID(status string) int64 {
	statusMap := map[string]int64{
		"draft":     1,
		"pending":   2,
		"active":    3,
		"completed": 4,
		"rejected":  5,
		"cancelled": 6,
	}

	statusID := statusMap[status]
	return statusID
}
