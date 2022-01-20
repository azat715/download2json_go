package logger

import (
	"log"
	"os"
)

// GeneralLogger exported
var GeneralLogger *log.Logger

// ErrorLogger exported
var ErrorLogger *log.Logger

func init() {
	GeneralLogger = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)
}
