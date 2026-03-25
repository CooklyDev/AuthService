package internal

import (
	"fmt"
	"os"
	"time"
)

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

type PostgresConfig struct {
	Host     string
	Port     uint16
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresConfig() *PostgresConfig {
	var port uint16
	portStr := LookupEnvRequired("DB_PORT")
	_, err := fmt.Sscanf(portStr, "%d", &port)
	if err != nil {
		panic("invalid DB_PORT value: " + portStr)
	}
	return &PostgresConfig{
		Host:     LookupEnvRequired("DB_HOST"),
		Port:     port,
		User:     LookupEnvRequired("DB_USER"),
		Password: LookupEnvRequired("DB_PASSWORD"),
		DBName:   LookupEnvRequired("DB_NAME"),
		SSLMode:  LookupEnvRequired("DB_SSL_MODE"),
	}
}

type RedisConfig struct {
	Host     string
	Port     uint16
	Password string
}

func NewRedisConfig() *RedisConfig {
	password, exists := LookupEnvOptional("REDIS_PASSWORD")
	if !exists {
		password = ""
	}
	var port uint16
	portStr := LookupEnvRequired("REDIS_PORT")
	_, err := fmt.Sscanf(portStr, "%d", &port)
	if err != nil {
		panic("invalid REDIS_PORT value: " + portStr)
	}
	return &RedisConfig{
		Host:     LookupEnvRequired("REDIS_HOST"),
		Port:     port,
		Password: password,
	}
}

type AppConfig struct {
	AppPort            string
	SessionTTL         time.Duration
	SessionPrefix      string
	PublisherInterval  time.Duration
	PublisherBatchSize uint64
}

func convertSessionTTL(sessionTTLStr string) time.Duration {
	sessionTTL, err := time.ParseDuration(sessionTTLStr)
	if err != nil {
		panic("invalid SESSION_TTL value: " + sessionTTLStr)
	}
	return sessionTTL
}

func convertDurationOrDefault(value string, defaultValue time.Duration, envName string) time.Duration {
	if value == "" {
		return defaultValue
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		panic("invalid " + envName + " value: " + value)
	}

	return duration
}

func convertUint64OrDefault(value string, defaultValue uint64, envName string) uint64 {
	if value == "" {
		return defaultValue
	}

	var parsed uint64
	_, err := fmt.Sscanf(value, "%d", &parsed)
	if err != nil {
		panic("invalid " + envName + " value: " + value)
	}

	return parsed
}

func NewAppConfig() *AppConfig {
	publisherInterval, _ := LookupEnvOptional("PUBLISHER_INTERVAL")
	publisherBatchSize, _ := LookupEnvOptional("PUBLISHER_BATCH_SIZE")

	return &AppConfig{
		SessionTTL:         convertSessionTTL(LookupEnvRequired("SESSION_TTL")),
		AppPort:            LookupEnvRequired("APP_PORT"),
		SessionPrefix:      LookupEnvRequired("SESSION_PREFIX"),
		PublisherInterval:  convertDurationOrDefault(publisherInterval, time.Second, "PUBLISHER_INTERVAL"),
		PublisherBatchSize: convertUint64OrDefault(publisherBatchSize, 100, "PUBLISHER_BATCH_SIZE"),
	}
}
