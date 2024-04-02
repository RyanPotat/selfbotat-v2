package Log

import (
	"log"
	"os"
)

var (
	info *log.Logger
	warning *log.Logger
	debug *log.Logger
	errors *log.Logger
)

func init() {
	info = log.New(os.Stdout, "INFO: ", log.Ldate | log.Ltime)
	warning = log.New(os.Stdout, "WARNING: ", log.Ldate | log.Ltime)
	debug = log.New(os.Stdout, "DEBUG: ", log.Ldate | log.Ltime)
	errors = log.New(os.Stdout, "ERROR: ", log.Ldate | log.Ltime)
}

func Info(input ...any) {
	info.Println(input ...)
}

func Warning(input ...any) {
	warning.Println(input ...)
}

func Debug(input ...any) {
	debug.Println(input ...)
}

func Error(input ...any) {
	errors.Println(input ...)
}

func Infof(format string, v ...interface{}) {
	info.Printf(format, v...)
}

func Warningf(format string, v ...interface{}) {
	warning.Printf(format, v...)
}

func Debugf(format string, v ...interface{}) {
	debug.Printf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	errors.Printf(format, v...)
}