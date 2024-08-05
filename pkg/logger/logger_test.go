package logger

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.CreateTemp("", "logger_test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir.Name())

	// Set the log directory to the temporary directory
	logDir := filepath.Join(tmpDir.Name(), "logs")
	os.MkdirAll(logDir, 0755)

	// Set the logToFileAndTerminal flag to true
	logToFileAndTerminal := true

	// Call the Init function
	InitLogger(logToFileAndTerminal)

	// Check if the log files are created in the correct directory
	infoLogPath := filepath.Join(logDir, "info.log")
	errorLogPath := filepath.Join(logDir, "error.log")

	if _, err := os.Stat(infoLogPath); os.IsNotExist(err) {
		t.Errorf("Info log file not created: %v", err)
	}

	if _, err := os.Stat(errorLogPath); os.IsNotExist(err) {
		t.Errorf("Error log file not created: %v", err)
	}

	// Check if the loggers are initialized correctly
	if InfoLogger == nil {
		t.Error("InfoLogger not initialized")
	}

	if ErrorLogger == nil {
		t.Error("ErrorLogger not initialized")
	}

	// Check if the loggers are writing to the correct writers
	var infoBuf bytes.Buffer
	var errorBuf bytes.Buffer

	InfoLogger.SetOutput(&infoBuf)
	ErrorLogger.SetOutput(&errorBuf)

	InfoLogger.Print("Test Info Log")
	ErrorLogger.Print("Test Error Log")

	infoOutput := infoBuf.String()
	errorOutput := errorBuf.String()

	if infoOutput == "" {
		t.Error("InfoLogger not writing to the correct writer")
	}

	if errorOutput == "" {
		t.Error("ErrorLogger not writing to the correct writer")
	}

	// Check if the loggers are using the correct log prefix
	currentDate := time.Now().Format("2006/01/02")
	expectedInfoPrefix := "INFO: " + currentDate
	expectedErrorPrefix := "ERROR: " + currentDate

	if infoOutput[:len(expectedInfoPrefix)] != expectedInfoPrefix {
		t.Errorf("InfoLogger prefix incorrect: expected %q, got %q", expectedInfoPrefix, infoOutput[:len(expectedInfoPrefix)])
	}

	if errorOutput[:len(expectedErrorPrefix)] != expectedErrorPrefix {
		t.Errorf("ErrorLogger prefix incorrect: expected %q, got %q", expectedErrorPrefix, errorOutput[:len(expectedErrorPrefix)])
	}
}
