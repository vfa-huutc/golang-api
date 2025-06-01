package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// InitValidator initializes the validator engine and registers custom validation rules.
// This function is called during the application startup to ensure that
func InitValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("valid_birthday", ValidateBirthday)
	}
}

// ValidateBirthday checks if the birthday is in a valid format and not a future date.
func ValidateBirthday(fl validator.FieldLevel) bool {
	birthdayStr := fl.Field().String()
	layout := "2006-01-02" // Format: YYYY-MM-DD

	// Parse the birthday to check the format
	parsedDate, err := time.Parse(layout, birthdayStr)
	if err != nil {
		return false // Invalid date format
	}

	// Check if the birthday is in the future
	if parsedDate.After(time.Now()) {
		return false // Invalid: birthday can't be in the future
	}

	return true // Valid birthday
}

func TranslateValidationErrors(err error) error {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		// Get the first validation error
		fe := ve[0]

		field := strings.ToLower(fe.Field())
		param := fe.Param()
		var msg string

		switch fe.Tag() {
		case "required":
			msg = fmt.Sprintf("%s is required", field)
		case "email":
			msg = fmt.Sprintf("%s must be a valid email address", field)
		case "url":
			msg = fmt.Sprintf("%s must be a valid URL", field)
		case "uuid":
			msg = fmt.Sprintf("%s must be a valid UUID", field)
		case "len":
			msg = fmt.Sprintf("%s must be exactly %s characters long", field, param)
		case "min":
			msg = fmt.Sprintf("%s must be at least %s characters long or numeric", field, param)
		case "max":
			msg = fmt.Sprintf("%s must be at most %s characters long or numeric", field, param)
		case "eq":
			msg = fmt.Sprintf("%s must be equal to %s", field, param)
		case "ne":
			msg = fmt.Sprintf("%s must not be equal to %s", field, param)
		case "lt":
			msg = fmt.Sprintf("%s must be less than %s", field, param)
		case "lte":
			msg = fmt.Sprintf("%s must be less than or equal to %s", field, param)
		case "gt":
			msg = fmt.Sprintf("%s must be greater than %s", field, param)
		case "gte":
			msg = fmt.Sprintf("%s must be greater than or equal to %s", field, param)
		case "oneof":
			msg = fmt.Sprintf("%s must be one of [%s]", field, param)
		case "contains":
			msg = fmt.Sprintf("%s must contain '%s'", field, param)
		case "excludes":
			msg = fmt.Sprintf("%s must not contain '%s'", field, param)
		case "startswith":
			msg = fmt.Sprintf("%s must start with '%s'", field, param)
		case "endswith":
			msg = fmt.Sprintf("%s must end with '%s'", field, param)
		case "ip":
			msg = fmt.Sprintf("%s must be a valid IP address", field)
		case "ipv4":
			msg = fmt.Sprintf("%s must be a valid IPv4 address", field)
		case "ipv6":
			msg = fmt.Sprintf("%s must be a valid IPv6 address", field)
		case "datetime":
			msg = fmt.Sprintf("%s must be a valid datetime (format: %s)", field, param)
		case "numeric":
			msg = fmt.Sprintf("%s must be a numeric value", field)
		case "boolean":
			msg = fmt.Sprintf("%s must be a boolean value", field)
		case "alpha":
			msg = fmt.Sprintf("%s must contain only letters", field)
		case "alphanum":
			msg = fmt.Sprintf("%s must contain only letters and numbers", field)
		case "alphanumunicode":
			msg = fmt.Sprintf("%s must contain only unicode letters and numbers", field)
		case "ascii":
			msg = fmt.Sprintf("%s must contain only ASCII characters", field)
		case "printascii":
			msg = fmt.Sprintf("%s must contain only printable ASCII characters", field)
		case "base64":
			msg = fmt.Sprintf("%s must be a valid base64 string", field)
		case "containsany":
			msg = fmt.Sprintf("%s must contain at least one of the characters in '%s'", field, param)
		case "excludesall":
			msg = fmt.Sprintf("%s must not contain any of the characters in '%s'", field, param)
		case "excludesrune":
			msg = fmt.Sprintf("%s must not contain the rune '%s'", field, param)
		case "isdefault":
			msg = fmt.Sprintf("%s must be the default value", field)
		case "unique":
			msg = fmt.Sprintf("%s must contain unique values", field)
		case "valid_birthday":
			msg = fmt.Sprintf("%s must be a valid date (YYYY-MM-DD) and not in the future", field)

		default:
			msg = fmt.Sprintf("%s is invalid", field)
		}

		return errors.New(msg)
	}
	return err
}
