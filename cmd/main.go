package main

import (
	"os"

	"github.com/stbit/gopack/pkg/manager"
)

func main() {
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	m, err := manager.NewManager(path)
	if err != nil {
		panic(err)
	}

	if err = m.Run(); err != nil {
		panic(err)
	}
}
