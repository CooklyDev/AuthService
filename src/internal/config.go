package internal

import "os"

func LookupEnvOptional(key string) (string, bool) {
	return os.LookupEnv(key)
}

func LookupEnvRequired(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		panic("required environment variable " + key + " is not set")
	}
	return value
}
