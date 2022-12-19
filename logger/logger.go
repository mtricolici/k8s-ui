package logger

import (
	"log"
	"os"
)

const (
	log_file = "app.log"
)

var (
	Log     *log.Logger
	logFile *os.File
)

func Init() {
	logFile, err := os.Create(log_file)
	if err != nil {
		log.Fatalf("Error creating file '%s'", log_file)
	}

	Log = log.New(logFile, "", log.LstdFlags|log.Lshortfile)
}

func Close() {
	logFile.Close()
}
