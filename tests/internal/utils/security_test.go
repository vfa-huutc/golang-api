package utils

import (
	"reflect"
	"testing"

	"github.com/vfa-khuongdv/golang-cms/internal/utils"
)

func TestCensorSensitiveData(t *testing.T) {
	maskFields := []string{"password", "apiKey"}

	t.Run("Nil input", func(t *testing.T) {
		var input any
		result := utils.CensorSensitiveData(input, maskFields)
		if result != nil {
			t.Errorf("Expected nil, got %v", result)
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
}
