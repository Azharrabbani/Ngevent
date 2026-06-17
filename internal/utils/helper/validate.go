package helper

import (
	"regexp"
	"strings"
)

func IsValidCategory(name string) bool {
	name = strings.TrimSpace(name)

	regex := regexp.MustCompile(`^[\p{L}\p{N}\s&+./-]+$`)

	return len(name) > 0 && regex.MatchString(name)
}

func NormalizeCategory(name string) string {
	name = strings.TrimSpace(name)

	name = regexp.MustCompile(`\d+$`).ReplaceAllString(name, "")
	name = strings.Join(strings.Fields(name), " ")

	return name
}
