package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Logger interface defines methods for various logging levels.
type Logger interface {
	LogInfo(msg ...interface{})
	LogError(msg ...interface{})
	LogFatal(msg ...interface{})
	Close() error
}

// FileLogger implements the Logger interface with file and stdout logging.
type FileLogger struct {
	logFile *os.File
	logger  *log.Logger
}

// NewLogger initializes a new FileLogger with the specified log file path.
// It creates necessary directories if they do not exist, and configures
// logging to both the specified file and stdout.
func NewLogger(logFilePath string) (Logger, error) {
	// Ensure the log directory exists
	logDir := filepath.Dir(logFilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open the log file
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening log file %s, using stdout only: %v", logFilePath, err)
		log.SetOutput(os.Stdout)
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	// Set up multi-writer to log to both stdout and the log file
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger := log.New(multiWriter, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

	return &FileLogger{
		logFile: logFile,
		logger:  logger,
	}, nil
}

// LogInfo logs an informational message.
func (f *FileLogger) LogInfo(msg ...interface{}) {
	f.logger.SetPrefix("INFO: ")
	f.logger.Println(msg...)
}

// LogError logs an error message.
func (f *FileLogger) LogError(msg ...interface{}) {
	f.logger.SetPrefix("ERROR: ")
	f.logger.Println(msg...)
}

// LogFatal logs a fatal error message and exits the application.
func (f *FileLogger) LogFatal(msg ...interface{}) {
	f.logger.SetPrefix("FATAL: ")
	f.logger.Fatalln(msg...)
}

// Close closes the log file.
func (f *FileLogger) Close() error {
	if f.logFile != nil {
		return f.logFile.Close()
	}
	return nil
}
