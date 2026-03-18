package internal

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/CooklyDev/AuthService/internal/domain"
)

const (
	colorReset  = "\033[0m"
	colorCyan   = "\033[36m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
)

type ConsoleLogger struct {
	out io.Writer
}

var _ domain.Logger = (*ConsoleLogger)(nil)

// NewConsoleLogger returns a logger that writes to stderr, which is the
// conventional destination for log output to avoid mixing with program output.
func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{out: os.Stderr}
}

func (l *ConsoleLogger) log(level, color, msg string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	if _, err := fmt.Fprintf(l.out, "%s%s [%-5s]%s %s\n", color, timestamp, level, colorReset, msg); err != nil {
		return
	}
}

func (l *ConsoleLogger) Debug(msg string) {
	l.log("DEBUG", colorCyan, msg)
}

func (l *ConsoleLogger) Info(msg string) {
	l.log("INFO", colorGreen, msg)
}

func (l *ConsoleLogger) Warn(msg string) {
	l.log("WARN", colorYellow, msg)
}

func (l *ConsoleLogger) Error(msg string) {
	l.log("ERROR", colorRed, msg)
}
