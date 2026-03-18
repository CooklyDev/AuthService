package internal

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
)

// ansiEscape matches ANSI escape sequences so they can be stripped from output.
var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(s string) string {
	return ansiEscape.ReplaceAllString(s, "")
}

func newTestConsoleLogger() (*ConsoleLogger, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	return &ConsoleLogger{out: buf}, buf
}

func TestConsoleLogger_Debug(t *testing.T) {
	// Arrange
	logger, buf := newTestConsoleLogger()

	// Act
	logger.Debug("debug message")

	// Assert
	output := buf.String()
	if !strings.Contains(output, "DEBUG") {
		t.Fatalf("expected DEBUG in output, got %q", output)
	}
	if !strings.Contains(output, "debug message") {
		t.Fatalf("expected message in output, got %q", output)
	}
}

func TestConsoleLogger_Info(t *testing.T) {
	// Arrange
	logger, buf := newTestConsoleLogger()

	// Act
	logger.Info("info message")

	// Assert
	output := buf.String()
	if !strings.Contains(output, "INFO") {
		t.Fatalf("expected INFO in output, got %q", output)
	}
	if !strings.Contains(output, "info message") {
		t.Fatalf("expected message in output, got %q", output)
	}
}

func TestConsoleLogger_Warn(t *testing.T) {
	// Arrange
	logger, buf := newTestConsoleLogger()

	// Act
	logger.Warn("warn message")

	// Assert
	output := buf.String()
	if !strings.Contains(output, "WARN") {
		t.Fatalf("expected WARN in output, got %q", output)
	}
	if !strings.Contains(output, "warn message") {
		t.Fatalf("expected message in output, got %q", output)
	}
}

func TestConsoleLogger_Error(t *testing.T) {
	// Arrange
	logger, buf := newTestConsoleLogger()

	// Act
	logger.Error("error message")

	// Assert
	output := buf.String()
	if !strings.Contains(output, "ERROR") {
		t.Fatalf("expected ERROR in output, got %q", output)
	}
	if !strings.Contains(output, "error message") {
		t.Fatalf("expected message in output, got %q", output)
	}
}

func TestConsoleLogger_OutputContainsTimestamp(t *testing.T) {
	// Arrange
	logger, buf := newTestConsoleLogger()

	// Act
	logger.Info("timestamp test")

	// Assert
	stripped := stripANSI(buf.String())
	// Timestamp format: "YYYY-MM-DD HH:MM:SS"
	if len(stripped) < 19 {
		t.Fatalf("expected output to contain a timestamp, got %q", stripped)
	}
	if stripped[4] != '-' || stripped[7] != '-' || stripped[13] != ':' || stripped[16] != ':' {
		t.Fatalf("expected timestamp in YYYY-MM-DD HH:MM:SS format, got %q", stripped)
	}
}
