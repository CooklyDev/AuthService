package internal

import "github.com/CooklyDev/AuthService/internal/domain"

type NoopLogger struct{}

var _ domain.Logger = (*NoopLogger)(nil)

func NewNoopLogger() *NoopLogger {
	return &NoopLogger{}
}

func (logger *NoopLogger) Debug(string) {}

func (logger *NoopLogger) Info(string) {}

func (logger *NoopLogger) Warn(string) {}

func (logger *NoopLogger) Error(string) {}
