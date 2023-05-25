package logger

import (
	"log"
)

func Error(err error) {
	log.Printf("%s: %v", Red("error"), err)
}

func Fatal(err error) {
	log.Fatalf("%s: %v", Red("error"), err)
}
