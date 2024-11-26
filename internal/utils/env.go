package utils

import (
	"os"
	"strconv"
)

// GetEnv retrieves a string value from the environment with a fallback default value
func GetEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// GetEnvAsInt retrieves an integer value from the environment with a fallback default value
func GetEnvAsInt(key string, defaultValue int) int {
	valueStr := GetEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
