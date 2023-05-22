package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/stbit/gopack/pkg/manager"
	"github.com/stbit/gopack/pkg/manager/execute"
)

func main() {
	var commandsExec execute.CommandsFlag

	watch := flag.Bool("w", false, "watch files changes")
	flag.Var(&commandsExec, "e", "execute commands after compile")
	flag.Parse()

	fmt.Println("v.0.0.1")

	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	m, err := manager.New(path)
	if err != nil {
		panic(err)
	}

	if err = m.Run(commandsExec); err != nil {
		panic(err)
	}

	if *watch {
		m.Watch(commandsExec)
	}
}
