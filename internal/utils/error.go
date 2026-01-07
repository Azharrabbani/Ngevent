package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func CustomErrorMessage(err validator.FieldError) ValidationError {
	field := ToSnakeCase(err.Field())
	param := ToWord(err.Param())

	var msg string
	switch err.Tag() {
	case "required":
		msg = fmt.Sprintf("%s is required.", field)
	case "oneof":
		msg = fmt.Sprintf("%s must be one of [%s].", field, param)
	default:
		msg = fmt.Sprintf("%s is invalid.", field)
	}

	return ValidationError{
		Field:   field,
		Message: msg,
	}
}

func GetValidationError(err error) []ValidationError {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make([]ValidationError, len(ve))
		for i, fe := range ve {
			out[i] = CustomErrorMessage(fe)
		}
		return out
	}

	return nil
}

func ToSnakeCase(str string) string {
	text := regexp.MustCompile("([a-z0-9])([A-Z])")
	result := text.ReplaceAllString(str, "${1}_${2}")

	return strings.ToLower(result)
}

func ToWord(str string) string {
	text := regexp.MustCompile("([a-z0-9])([A-Z])")
	result := text.ReplaceAllString(str, "${1} ${2}")

	return strings.ToLower(result)
}
