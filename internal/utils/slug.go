package utils

import (
	"regexp"
	"strings"
)

func CreateSlug(s string) string {
	lower := strings.ToLower(s)

	re := regexp.MustCompile(`(?m)[^a-zA-Z0-9]+`)
	result := re.ReplaceAllString(lower, "-")

	return strings.Trim(result, "-")
}
