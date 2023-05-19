package main

import (
	"log"
	"os"
	"time"

	"github.com/stbit/gopack/pkg/manager"
)

func main() {
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
