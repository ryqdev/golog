package golog

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

// TestLoggerInitialization checks that the logger initializes with the default level and output.
func TestLoggerInitialization(t *testing.T) {
	if GetLevel() != LevelInfo {
		t.Errorf("expected log level %v, got %v", LevelInfo, GetLevel())
	}
	if defaultLogger.w != os.Stderr {
		t.Error("expected default writer to be os.Stderr")
	}
}

// TestSetLevel checks the SetLevel method.
func TestSetLevel(t *testing.T) {
	SetLevel(LevelDebug)
	if GetLevel() != LevelDebug {
		t.Errorf("expected log level %v, got %v", LevelDebug, GetLevel())
	}
}

// TestInfoLogging checks that Info messages are correctly logged.
func TestInfoLogging(t *testing.T) {
	var buf bytes.Buffer
	defaultLogger.w = &buf
	SetLevel(LevelInfo)
	Info("test info message")

	expected := fmt.Sprintf("%s test info message \n", InfoLevel)
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

// TestDebugLogging checks that Debug messages are correctly logged.
func TestDebugLogging(t *testing.T) {
	var buf bytes.Buffer
	defaultLogger.w = &buf
	SetLevel(LevelDebug)
	Debug("test debug message")

	expected := fmt.Sprintf("%s test debug message \n", DebugLevel)
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

// TestErrorLogging checks that Error messages are correctly logged.
func TestErrorLogging(t *testing.T) {
	var buf bytes.Buffer
	defaultLogger.w = &buf
	SetLevel(LevelError)
	Error("test error message")

	expected := fmt.Sprintf("%s test error message \n", ErrorLevel)
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

// TestAddProcessor checks that custom processors are applied correctly.
func TestAddProcessor(t *testing.T) {
	var buf bytes.Buffer
	defaultLogger.w = &buf

	// Define a processor that adds a prefix to the log message
	prefixProcessor := func(format string, v ...any) (string, []any) {
		return "[PREFIX] " + format, v
	}
	AddProcessor(prefixProcessor)

	SetLevel(LevelInfo)
	Info("test info message")

	expected := fmt.Sprintf("%s [PREFIX] test info message \n", InfoLevel)
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}
