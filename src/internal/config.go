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
	AppPort       string
	SessionTTL    time.Duration
	SessionPrefix string
}

func convertSessionTTL(sessionTTLStr string) time.Duration {
	sessionTTL, err := time.ParseDuration(sessionTTLStr)
	if err != nil {
		panic("invalid SESSION_TTL value: " + sessionTTLStr)
	}
	return sessionTTL
}

func NewAppConfig() *AppConfig {
	return &AppConfig{
		SessionTTL:    convertSessionTTL(LookupEnvRequired("SESSION_TTL")),
		AppPort:       LookupEnvRequired("APP_PORT"),
		SessionPrefix: LookupEnvRequired("SESSION_PREFIX"),
	}
}
