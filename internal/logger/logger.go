package logger

import (
	"log"
	"os"
)

var logger *log.Logger

// Init initializes the logger by creating or opening the specified log file. It sets up the logger to write log messages with timestamps. It returns an error if there is an issue opening the file.
func Init(filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	logger = log.New(file, " ", log.LstdFlags)
	return nil
}

// Log writes a log message to the log file. If the logger is not initialized, it does nothing.
func Log(message string) {
	if logger != nil {
		logger.Println(message)
	}
}
