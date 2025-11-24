package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log *logrus.Logger

// InitLogger initializes the application logger with file rotation
func InitLogger() error {
	Log = logrus.New()

	// Get log directory path
	logDir, err := getLogDirectory()
	if err != nil {
		return fmt.Errorf("failed to get log directory: %w", err)
	}

	// Create log directory if it doesn't exist
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	logFile := filepath.Join(logDir, "app.log")

	// Set up log file with rotation
	fileWriter := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10, // megabytes
		MaxBackups: 5,
		MaxAge:     30, // days
		Compress:   true,
	}

	// Log to both file and stdout for development
	Log.SetOutput(fileWriter)
	
	// Set formatter to JSON for structured logging
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Set log level (can be changed to Debug for more verbose logging)
	Log.SetLevel(logrus.InfoLevel)

	Log.WithFields(logrus.Fields{
		"log_file": logFile,
	}).Info("Logger initialized successfully")

	return nil
}

// getLogDirectory returns the path to the log directory
func getLogDirectory() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".mlr-desktop", "logs"), nil
}

// RecoverFromPanic recovers from a panic and logs it with stack trace
func RecoverFromPanic(functionName string) {
	if r := recover(); r != nil {
		stackTrace := string(debug.Stack())
		
		// Only log if logger is initialized
		if Log != nil {
			Log.WithFields(logrus.Fields{
				"function":    functionName,
				"panic_value": r,
				"stack_trace": stackTrace,
			}).Error("Panic recovered")
		}
		
		// Also print to stderr for immediate visibility
		fmt.Fprintf(os.Stderr, "PANIC in %s: %v\n%s\n", functionName, r, stackTrace)
	}
}

// LogError is a helper to log errors with context
func LogError(err error, context string, fields logrus.Fields) {
	if err == nil || Log == nil {
		return
	}
	
	if fields == nil {
		fields = logrus.Fields{}
	}
	fields["context"] = context
	fields["error"] = err.Error()
	
	Log.WithFields(fields).Error("Error occurred")
}

// LogInfo is a helper to log info messages with context
func LogInfo(message string, fields logrus.Fields) {
	if Log == nil {
		return
	}
	if fields == nil {
		fields = logrus.Fields{}
	}
	Log.WithFields(fields).Info(message)
}

// LogWarn is a helper to log warning messages with context
func LogWarn(message string, fields logrus.Fields) {
	if Log == nil {
		return
	}
	if fields == nil {
		fields = logrus.Fields{}
	}
	Log.WithFields(fields).Warn(message)
}
