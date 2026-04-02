package helper

import (
	"fmt"
	"time"
)

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

func StrPointerIfNotEmpty(s string) *string {
	if s == "" {
		return nil
	}

	return &s
}

func ArrayIntToPointer(i []int) *[]int {
	if len(i) == 0 {
		return nil
	}

	return &i
}

func TimeToPointer(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}

	return &t
}

func TimePtrToUnix(t *time.Time) *int64 {
	if t == nil {
		return nil
	}
	val := t.Unix()
	return &val
}
