package config

import (
	"os"
	"strconv"
	"strings"
)

// LookupEnv is a generic type implementation to search env keys
func LookupEnv[T string | []string | int](name string, defaultValue T) T {
	value, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}

	var result any
	switch any(defaultValue).(type) {
	case string:
		result = value
	case []string:
		// Should be a comma separated list
		strs := strings.Split(value, ",")
		result = strs
	case int:
		i, _ := strconv.ParseInt(value, 10, 64)
		result = int(i)
	}

	return result.(T)
}
