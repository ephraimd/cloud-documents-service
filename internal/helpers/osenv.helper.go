package helpers

import (
	"os"
)

func GetOsEnvOrDefault(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func GetOsEnvOrPanic(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	panic("Environment variable " + key + " is not set")
}
