package utils_test

import (
	"testing"
	"time"

	"errors"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
)

type User struct {
	Birthday string `validate:"required,valid_birthday"`
}

type TestStruct struct {
	Field string `validate:"%s=%s"`
}

func TestValidateBirthday(t *testing.T) {
	validate := validator.New()
	validate.RegisterValidation("valid_birthday", utils.ValidateBirthday)

	tests := []struct {
		name     string
		birthday string
		wantErr  bool
	}{
		{
			name:     "Valid birthday",
			birthday: "2000-01-01",
			wantErr:  false,
		},
		{
			name:     "Invalid format",
			birthday: "01-01-2000",
			wantErr:  true,
		},
		{
			name:     "Future date",
			birthday: time.Now().AddDate(1, 0, 0).Format("2006-01-02"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := User{Birthday: tt.birthday}
			err := validate.Struct(u)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTranslateValidationErrors(t *testing.T) {
	validate := validator.New()

	testCases := []struct {
		name     string
		tag      string
		param    string
		value    any
		expected string
	}{
		{"required", "required", "", struct {
			Field string `validate:"required"`
		}{}, "field is required"},

		{"email", "email", "", struct {
			Field string `validate:"email"`
		}{Field: "invalid"}, "field must be a valid email address"},

		{"url", "url", "", struct {
			Field string `validate:"url"`
		}{Field: "not-a-url"}, "field must be a valid URL"},

		{"uuid", "uuid", "", struct {
			Field string `validate:"uuid"`
		}{Field: "invalid-uuid"}, "field must be a valid UUID"},

		{"len", "len", "5", struct {
			Field string `validate:"len=5"`
		}{Field: "abc"}, "field must be exactly 5 characters long"},

		{"min", "min", "5", struct {
			Field string `validate:"min=5"`
		}{Field: "abc"}, "field must be at least 5 characters long or numeric"},

		{"max", "max", "2", struct {
			Field string `validate:"max=2"`
		}{Field: "abc"}, "field must be at most 2 characters long or numeric"},

		{"eq", "eq", "admin", struct {
			Field string `validate:"eq=admin"`
		}{Field: "user"}, "field must be equal to admin"},

		{"ne", "ne", "admin", struct {
			Field string `validate:"ne=admin"`
		}{Field: "admin"}, "field must not be equal to admin"},

		{"lt", "lt", "10", struct {
			Field int `validate:"lt=10"`
		}{Field: 20}, "field must be less than 10"},

		{"lte", "lte", "10", struct {
			Field int `validate:"lte=10"`
		}{Field: 20}, "field must be less than or equal to 10"},

		{"gt", "gt", "10", struct {
			Field int `validate:"gt=10"`
		}{Field: 5}, "field must be greater than 10"},

		{"gte", "gte", "10", struct {
			Field int `validate:"gte=10"`
		}{Field: 5}, "field must be greater than or equal to 10"},

		{"oneof", "oneof", "admin user", struct {
			Field string `validate:"oneof=admin user"`
		}{Field: "guest"}, "field must be one of [admin user]"},

		{"contains", "contains", "@", struct {
			Field string `validate:"contains=@"`
		}{Field: "example.com"}, "field must contain '@'"},

		{"excludes", "excludes", "@", struct {
			Field string `validate:"excludes=@"`
		}{Field: "user@example.com"}, "field must not contain '@'"},

		{"startswith", "startswith", "abc", struct {
			Field string `validate:"startswith=abc"`
		}{Field: "xyz"}, "field must start with 'abc'"},

		{"endswith", "endswith", "xyz", struct {
			Field string `validate:"endswith=xyz"`
		}{Field: "abc"}, "field must end with 'xyz'"},

		{"ip", "ip", "", struct {
			Field string `validate:"ip"`
		}{Field: "not-an-ip"}, "field must be a valid IP address"},

		{"ipv4", "ipv4", "", struct {
			Field string `validate:"ipv4"`
		}{Field: "not-an-ip"}, "field must be a valid IPv4 address"},

		{"ipv6", "ipv6", "", struct {
			Field string `validate:"ipv6"`
		}{Field: "not-an-ip"}, "field must be a valid IPv6 address"},

		{"datetime", "datetime", "2006-01-02", struct {
			Field string `validate:"datetime=2006-01-02"`
		}{Field: "01-01-2023"}, "field must be a valid datetime (format: 2006-01-02)"},

		{"numeric", "numeric", "", struct {
			Field string `validate:"numeric"`
		}{Field: "abc"}, "field must be a numeric value"},

		{"boolean", "boolean", "", struct {
			Field string `validate:"boolean"`
		}{Field: "notabool"}, "field must be a boolean value"},

		{"alpha", "alpha", "", struct {
			Field string `validate:"alpha"`
		}{Field: "abc123"}, "field must contain only letters"},

		{"alphanum", "alphanum", "", struct {
			Field string `validate:"alphanum"`
		}{Field: "abc!"}, "field must contain only letters and numbers"},

		{"alphanumunicode", "alphanumunicode", "", struct {
			Field string `validate:"alphanumunicode"`
		}{Field: "abc漢字!"}, "field must contain only unicode letters and numbers"},

		{"ascii", "ascii", "", struct {
			Field string `validate:"ascii"`
		}{Field: "こんにちは"}, "field must contain only ASCII characters"},

		{"printascii", "printascii", "", struct {
			Field string `validate:"printascii"`
		}{Field: "\u0000"}, "field must contain only printable ASCII characters"},

		{"base64", "base64", "", struct {
			Field string `validate:"base64"`
		}{Field: "not_base64"}, "field must be a valid base64 string"},

		{"containsany", "containsany", "abc", struct {
			Field string `validate:"containsany=abc"`
		}{Field: "xyz"}, "field must contain at least one of the characters in 'abc'"},

		{"excludesall", "excludesall", "abc", struct {
			Field string `validate:"excludesall=abc"`
		}{Field: "cab"}, "field must not contain any of the characters in 'abc'"},

		{"excludesrune", "excludesrune", "あ", struct {
			Field string `validate:"excludesrune=あ"`
		}{Field: "あいう"}, "field must not contain the rune 'あ'"},

		{"isdefault", "isdefault", "", struct {
			Field string `validate:"isdefault"`
		}{Field: "notdefault"}, "field must be the default value"},

		{"unique", "unique", "", struct {
			Field []string `validate:"unique"`
		}{Field: []string{"a", "b", "a"}}, "field must contain unique values"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validate.Struct(tc.value)
			assert.Error(t, err)

			transErr := utils.TranslateValidationErrors(err)
			assert.EqualError(t, transErr, tc.expected)
		})
	}
}

func TestTranslateValidationErrors_ExtraCases(t *testing.T) {
	validate := validator.New()

	validate.RegisterValidation("valid_birthday", utils.ValidateBirthday)
	validate.RegisterValidation("unknown", func(fl validator.FieldLevel) bool {
		return false
	})

	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{
			name: "valid_birthday (future date)",
			input: struct {
				Field string `validate:"valid_birthday"`
			}{Field: "3000-01-01"},
			expected: "field must be a valid date (YYYY-MM-DD) and not in the future",
		},
		{
			name: "default case fallback with unknown tag",
			input: struct {
				Field string `validate:"unknown"` // Custom registered tag
			}{Field: "data"},
			expected: "field is invalid",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := validate.Struct(tc.input)
			assert.Error(t, err)

			result := utils.TranslateValidationErrors(err)
			assert.EqualError(t, result, tc.expected)
		})
	}
}

func TestTranslateValidationErrors_NonValidationError(t *testing.T) {
	plainErr := errors.New("some generic error")

	result := utils.TranslateValidationErrors(plainErr)

	assert.Error(t, result)
	assert.EqualError(t, result, "some generic error")
}

func TestInitValidator(t *testing.T) {
	// Initialize the validator and register custom validations
	utils.InitValidator()

	// Get the validator engine from gin binding
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		t.Fatal("Failed to get validator engine")
	}

	// Check if the "valid_birthday" validation function is registered
	err := v.Var("3000-01-01", "valid_birthday") // Future date - should fail validation
	if err == nil {
		t.Error("Expected validation error for future birthday, got nil")
	}

	err = v.Var("2000-01-01", "valid_birthday") // Valid date - should pass validation
	if err != nil {
		t.Errorf("Expected no validation error for valid birthday, got: %v", err)
	}
}
