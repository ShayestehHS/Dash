package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// GetString retrieves a string environment variable.
// Returns empty string if the variable is not set.
func GetString(key string) string {
	return os.Getenv(key)
}

// GetBool retrieves a boolean environment variable.
// Returns the boolean value and nil error if the variable is set and valid.
// Returns false and an error if the variable is not set or cannot be converted to boolean.
func GetBool(key string) (bool, error) {
	value := os.Getenv(key)
	if value == "" {
		return false, fmt.Errorf("environment variable %s is not set", key)
	}

	upperValue := strings.ToUpper(value)
	switch upperValue {
	case "TRUE", "1", "YES", "ON":
		return true, nil
	case "FALSE", "0", "NO", "OFF":
		return false, nil
	default:
		return false, fmt.Errorf("failed to parse environment variable %s=%s as bool", key, value)
	}
}

// GetInt retrieves an integer environment variable.
// Returns the integer value and nil error if the variable is set and valid.
// Returns 0 and an error if the variable is not set or cannot be converted to integer.
func GetInt(key string) (int, error) {
	value := os.Getenv(key)
	if value == "" {
		return 0, fmt.Errorf("environment variable %s is not set", key)
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("failed to parse environment variable %s=%s as int: %w", key, value, err)
	}

	return intValue, nil
}

// GetInt64 retrieves an int64 environment variable.
// Returns the int64 value and nil error if the variable is set and valid.
// Returns 0 and an error if the variable is not set or cannot be converted to int64.
func GetInt64(key string) (int64, error) {
	value := os.Getenv(key)
	if value == "" {
		return 0, fmt.Errorf("environment variable %s is not set", key)
	}

	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse environment variable %s=%s as int64: %w", key, value, err)
	}

	return intValue, nil
}


