package utils_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
)

type secretStringer struct {
	secret string
}

func (s secretStringer) String() string {
	return s.secret
}

func TestCensorSensitiveData(t *testing.T) {
	maskFields := []string{"password", "apiKey"}

	t.Run("Nil value input", func(t *testing.T) {
		type TestInput struct {
			Name     string
			Password *string
			Deeps    struct {
				Password *string
			}
			DeepsPtr *struct {
				Password *string
			}
		}

		passwordPtr := "12345"
		passwordPtrMasked := "1***5"

		input := TestInput{
			Name: "test", Password: nil,
			Deeps: struct {
				Password *string
			}{Password: nil},
			DeepsPtr: &struct {
				Password *string
			}{Password: &passwordPtr}}

		expected := TestInput{Name: "test", Password: nil, Deeps: struct {
			Password *string
		}{Password: nil}, DeepsPtr: &struct {
			Password *string
		}{Password: &passwordPtrMasked}}

		result := utils.CensorSensitiveData(input, maskFields)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}

		// nil value input without maskFields
		result = utils.CensorSensitiveData(nil, maskFields)
		assert.Nil(t, result, "Expected nil result for nil input without maskFields")

	})

	t.Run("Nil input", func(t *testing.T) {
		var input interface{} = nil
		result := utils.CensorSensitiveData(input, maskFields)
		if result != nil {
			t.Errorf("Expected nil, got %v", result)
		}
	})
	t.Run("Nil value input without maskFields", func(t *testing.T) {
		type TestInput struct {
			Name     string
			DeepsPtr *struct {
				Name *string
			}
		}
		// input
		input := TestInput{
			Name:     "test",
			DeepsPtr: &struct{ Name *string }{Name: nil},
		}
		// expected
		expected := TestInput{
			Name:     "test",
			DeepsPtr: &struct{ Name *string }{Name: nil},
		}
		// call function
		result := utils.CensorSensitiveData(input, nil)
		// check result
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}

	})

	t.Run("Map with sensitive keys", func(t *testing.T) {
		input := map[string]string{"password": "secret", "username": "user"}
		expected := map[string]string{"password": "s****t", "username": "user"}
		result := utils.CensorSensitiveData(input, maskFields)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Struct with sensitive fields", func(t *testing.T) {
		type User struct {
			Password string
			Username string
		}
		input := User{Password: "secret", Username: "user"}
		expected := User{Password: "s****t", Username: "user"}
		result := utils.CensorSensitiveData(input, maskFields).(User)
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Pointer to struct", func(t *testing.T) {
		type User struct {
			Password string
			Username string
		}
		input := &User{Password: "secret", Username: "user"}
		expected := User{Password: "s****t", Username: "user"}
		result := utils.CensorSensitiveData(input, maskFields)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Nested structures", func(t *testing.T) {
		type Profile struct {
			APIKey  string
			Details map[string]string
		}
		type User struct {
			Password string
			Profile  Profile
		}
		input := User{
			Password: "mypassword",
			Profile: Profile{
				APIKey: "12345",
				Details: map[string]string{
					"username": "user",
					"email":    "user@example.com",
				},
			},
		}
		expected := User{
			Password: "m********d",
			Profile: Profile{
				APIKey: "1***5",
				Details: map[string]string{
					"username": "user",
					"email":    "user@example.com",
				},
			},
		}
		result := utils.CensorSensitiveData(input, maskFields).(User)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Slice of structs", func(t *testing.T) {
		type User struct {
			Password string
			Username string
		}
		input := []User{
			{Password: "secret", Username: "user"},
		}
		expected := []User{
			{Password: "s****t", Username: "user"},
		}
		result := utils.CensorSensitiveData(input, maskFields).([]User)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Ptr value is nil", func(t *testing.T) {
		type testStruct struct {
			Field1 string
			Field2 string
		}

		var ptr *testStruct = nil

		result := utils.CensorSensitiveData(ptr, []string{"Field1"})
		if result != nil {
			t.Errorf("Expected nil result for nil pointer input, got: %#v", result)
		}
	})

	t.Run("Mask string fields in struct", func(t *testing.T) {
		type User struct {
			Password string
			Username string
		}

		maskFields := []string{"Password"}

		input := User{Password: "1", Username: "user"}
		expected := User{
			Password: "*", // Assuming the mask function replaces the password with a single asterisk
			Username: "user",
		}

		result := utils.CensorSensitiveData(input, maskFields).(User)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}
	})

	t.Run("Default case with unsupported type", func(t *testing.T) {
		maskFields := []string{"Password"}

		tests := []struct {
			name  string
			input any
			want  any
		}{
			{
				name:  "int input returns same int",
				input: 42,
				want:  42,
			},
			{
				name:  "float input returns same float",
				input: 3.14,
				want:  3.14,
			},
			{
				name:  "bool input returns same bool",
				input: true,
				want:  true,
			},
			{
				name:  "complex input returns same complex",
				input: complex(1, 2),
				want:  complex(1, 2),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := utils.CensorSensitiveData(tt.input, maskFields)
				if got != tt.want {
					t.Errorf("CensorSensitiveData() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("Test with value is type stringer", func(t *testing.T) {

		type testStruct struct {
			Secret string
			Other  string
		}

		maskFields := []string{"Secret"}

		input := testStruct{
			Secret: secretStringer{secret: "verysecret"}.String(), // truyền chuỗi từ Stringer
			Other:  "public",
		}

		got := utils.CensorSensitiveData(input, maskFields)

		result, ok := got.(testStruct)
		if !ok {
			t.Fatalf("expected testStruct, got %T", got)
		}

		expectedMasked := "v********t"
		if result.Secret != expectedMasked {
			t.Errorf("expected masked Secret %q, got %q", expectedMasked, result.Secret)
		}

		if result.Other != input.Other {
			t.Errorf("expected Other to be unchanged: got %q", result.Other)
		}

	})

	t.Run("Test with []byte input", func(t *testing.T) {
		type byteStruct struct {
			Secret []byte
			Other  string
		}

		maskFields := []string{"Secret"}

		input := byteStruct{
			Secret: []byte("verysecret"),
			Other:  "public",
		}
		expected := byteStruct{
			Secret: []byte("v********t"),
			Other:  "public",
		}

		result := utils.CensorSensitiveData(input, maskFields).(byteStruct)

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

}
