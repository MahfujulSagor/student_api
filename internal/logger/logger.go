package logger

import (
	"github/mahfujulsagor/student_api/internal/config"
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	Debug *log.Logger
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
)

func Init(cfg *config.Config, env string) {
	//? Determine log file path
	logFilePath := cfg.LoggingConfig.File
	if logFilePath == "" {
		logFilePath = "logs/app.log"
	}

	//? Ensure log directory exists
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatal("Failed to create logs directory:", err)
	}

	logFile, err := os.OpenFile(filepath.Join("logs", "app.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	//? In development, log to both console and file
	var writer io.Writer
	if env == "development" {
		writer = io.MultiWriter(os.Stdout, logFile)
	} else {
		writer = logFile
	}

	//? Define loggers with standard flags
	flags := log.Ldate | log.Ltime | log.Lshortfile
	Debug = log.New(writer, "DEBUG: ", flags)
	Info = log.New(writer, "INFO: ", flags)
	Warn = log.New(writer, "WARN: ", flags)
	Error = log.New(writer, "ERROR: ", flags)
}
