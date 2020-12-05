package util

import (
	"io"
	"log"
	"os"
)

func Logging(logFile string) {
	logfile, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Fatalf("logfile="+logFile, err)
	}
	log.SetOutput(io.MultiWriter(logfile, os.Stdout))
}
