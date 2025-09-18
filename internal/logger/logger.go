// ! Optimize for production
package logger

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

func Init() {
	if err := os.MkdirAll("logs", 0755); err != nil {
		log.Fatal("Failed to create logs directory:", err)
	}

	logFile, err := os.OpenFile(filepath.Join("logs", "app.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	//? Define loggers
	mw := io.MultiWriter(os.Stdout, logFile)
	Info = log.New(mw, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
