package helper

import "fmt"

func StringValue(s *string) string {
	if s == nil {
		return ""
	}

	return *s
}

func BoolValue(b *bool) string {
	if b == nil {
		return ""
	}

	return fmt.Sprintf("%t", *b)
}
