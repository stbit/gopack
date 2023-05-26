package main

import (
	"flag"
	"os"

	"github.com/stbit/gopack"
	"github.com/stbit/gopack/pkg/plugins"
	"github.com/stbit/gopack/pkg/plugins/syncerr"
)

func main() {
	watch := flag.Bool("w", false, "watch files")
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	flag.Parse()

	g := gopack.New(&gopack.Config{
		RootPath: path,
		Watch:    *watch,
		Commands: []gopack.Command{
			{Name: "main", Exec: "go run ./dist/example/cmd/main.go"},
		},
		Plugins: []plugins.PluginRegister{
			&syncerr.SyncErrPlugin{},
		},
	})

	if err := g.Run(); err != nil {
		panic(err)
	}
}
