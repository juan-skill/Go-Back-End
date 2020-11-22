package env

import (
	"os"
	"strconv"
	"strings"
)

// GetString gets the env var as a string
func GetString(varName string, defaultValue string) string {
	val, _ := os.LookupEnv(varName)
	if val == "" {
		return defaultValue
	}

	return val
}

// GetBool gets the env var as a boolean
func GetBool(varName string, defaultValue bool) bool {
	val, ok := os.LookupEnv(varName)
	if !ok {
		return defaultValue
	}

	switch val {
	case "1", "t", "T", "true", "TRUE", "True", "yes", "Yes", "YES":
		return true
	case "0", "f", "F", "false", "FALSE", "False", "no", "No", "NO":
		return false
	}

	return defaultValue
}

// GetInt64 gets the env var as an int
func GetInt64(varName string, defaultValue int64) int64 {
	val, ok := os.LookupEnv(varName)
	if !ok {
		return defaultValue
	}

	iVal, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return defaultValue
	}

	return iVal
}

// GetFloat64 gets the env var a float
func GetFloat64(varName string, defaultValue float64) float64 {
	val, ok := os.LookupEnv(varName)
	if !ok {
		return defaultValue
	}

	fVal, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return defaultValue
	}

	return fVal
}

// GetStringArray gets the env var a []string given a separator
func GetStringArray(varName string, sep string, defaultValue []string) []string {
	val, ok := os.LookupEnv(varName)
	if !ok {
		return defaultValue
	}

	return strings.Split(val, sep)
}
