package domain

import "strings"

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

func MaskEmail(email string) string {
	at := strings.Index(email, "@")
	if at <= 1 {
		return "***"
	}

	return email[:1] + "***" + email[at:]
}
