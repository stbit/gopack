package main

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/stbit/gopack/pkg/manager"
)

type (
	Suka = int
	Test struct{}

	FunS func(int) int
)

var q Suka = reflect.Zero(reflect.TypeOf((*Suka)(nil)).Elem()).Interface().(Suka)

func main() {
	v := reflect.Zero(reflect.TypeOf((*FunS)(nil)).Elem()).Interface().(FunS)

	fmt.Println(v, q)

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	start := time.Now()
	m, err := manager.NewManager(path) // + string(os.PathSeparator) + "example")
	if err != nil {
		panic(err)
	}
	startRun := time.Now()
	if err = m.Run(); err != nil {
		panic(err)
	}
	elapsed := time.Since(start)
	endRun := time.Since(startRun)
	log.Printf("Proccess %s", endRun)
	log.Printf("Finish %s", elapsed)
}
