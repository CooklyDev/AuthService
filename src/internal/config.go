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
	portStr := LookupEnvRequired("POSTGRES_PORT")
	_, err := fmt.Sscanf(portStr, "%d", &port)
	if err != nil {
		panic("invalid POSTGRES_PORT value: " + portStr)
	}
	return &PostgresConfig{
		Host:     LookupEnvRequired("POSTGRES_HOST"),
		Port:     port,
		User:     LookupEnvRequired("POSTGRES_USER"),
		Password: LookupEnvRequired("POSTGRES_PASSWORD"),
		DBName:   LookupEnvRequired("POSTGRES_DB"),
		SSLMode:  LookupEnvRequired("POSTGRES_SSL_MODE"),
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
	var sessionTTL time.Duration
	_, err := fmt.Sscanf(sessionTTLStr, "%d", &sessionTTL)
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
