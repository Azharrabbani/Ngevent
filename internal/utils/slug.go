package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

func CreateSlug(s string) string {
	lower := strings.ToLower(s)

	re := regexp.MustCompile(`(?m)[^a-zA-Z0-9]+`)
	result := re.ReplaceAllString(lower, "-")

	return strings.Trim(result, "-")
}

func GenerateEventSlug(name string) string {
	return fmt.Sprintf(
		"%s-%s",
		CreateSlug(name),
		strings.ToLower(uuid.New().String()[:6]),
	)
}
