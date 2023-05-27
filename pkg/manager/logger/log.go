package logger

import (
	"fmt"
	"log"
)

func unwrapError(err error) []error {
	switch x := err.(type) {
	case interface{ Unwrap() []error }:
		return x.Unwrap()
	case error:
		return []error{x}
	default:
		panic(fmt.Errorf("unkown log error"))
	}
}

func Error(err error) {
	for _, e := range unwrapError(err) {
		log.Printf("%s: %v", Red("error"), e)
	}
}

func Fatal(err error) {
	for _, e := range unwrapError(err) {
		log.Fatalf("%s: %v", Red("error"), e)
	}
}
