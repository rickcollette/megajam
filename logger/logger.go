package logger

import (
	"log"
	"os"
	"path/filepath"
)

var (
	Logger *log.Logger // Global logger
)

// InitLogger initializes the logger and makes it available globally.
func InitLogger(logPath string) (*os.File, error) {
	logFolder := filepath.Dir(logPath)

	// Ensure the folder exists
	err := os.MkdirAll(logFolder, os.ModePerm)
	if err != nil {
		return nil, err
	}

	// Open or create the log file
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// Initialize the global logger
	Logger = log.New(file, "", log.Ldate|log.Ltime|log.Lshortfile)
	Logger.Println("Logger initialized")
	return file, nil
}
