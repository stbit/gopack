package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/stbit/gopack/pkg/manager"
)

func main() {
	fmt.Println("t1")
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	start := time.Now()
	m, err := manager.New(path)
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
