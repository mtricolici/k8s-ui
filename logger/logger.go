package logger

import (
	"log"
	"os"
	"time"
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

func LogExecutedTime(name string) func() {
	start := time.Now()
	return func() {
		Log.Printf("%s - took %v", name, time.Since(start))
	}
}
