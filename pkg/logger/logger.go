package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

/*
Init initializes the logger with the given logToFileAndTerminal flag.
If logToFileAndTerminal is true, the logger will log to both the terminal and a file.
If logToFileAndTerminal is false, the logger will only log to the file.
*/
func InitLogger(logToFileAndTerminal bool) {
	if InfoLogger != nil && ErrorLogger != nil {
		// Logger already initialized
		return
	}

	logDir := filepath.Join("logs", time.Now().Format("2006-01-02"))
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatal("Failed to create log directory:", err)
	}

	// Open log files for today's date
	infoLogFile, err := os.OpenFile(filepath.Join(logDir, "info.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open info log file:", err)
	}

	errorLogFile, err := os.OpenFile(filepath.Join(logDir, "error.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open error log file:", err)
	}

	// Create multi-writer to write to both file and terminal if logToFileAndTerminal is true
	var infoWriter io.Writer = infoLogFile
	var errorWriter io.Writer = errorLogFile

	if logToFileAndTerminal {
		infoWriter = io.MultiWriter(os.Stdout, infoLogFile)
		errorWriter = io.MultiWriter(os.Stderr, errorLogFile)
	}

	// Initialize loggers
	InfoLogger = log.New(infoWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(errorWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

}

func CleanLogs() {
	// logDir := filepath.Join("logs", time.Now().Format("2006-01-02"))
	if err := os.RemoveAll("logs"); err != nil {
		log.Fatal("Failed to remove log directory:", err)
	}
}
