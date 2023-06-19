package main

import (
	"flag"
	"os"

	"github.com/iancoleman/strcase"
	"github.com/stbit/gopack"
	"github.com/stbit/gopack/plugins"
	"github.com/stbit/gopack/plugins/jsontag"
	"github.com/stbit/gopack/plugins/livereload"
	"github.com/stbit/gopack/plugins/syncerr"
)

func main() {
	commands := []gopack.Command{}
	watch := flag.Bool("w", false, "watch files")
	path, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	flag.Parse()

	if *watch {
		commands = []gopack.Command{
			{Name: "main", Exec: "go run ./dist/example/cmd/main.go"},
		}
	}

	g := gopack.New(&gopack.Config{
		RootPath: path,
		Watch:    *watch,
		Commands: commands,
		Plugins: []plugins.PluginRegister{
			syncerr.New(),
			jsontag.New(func(tag *jsontag.Tag) (name string, options []string) {
				if tag.JsonName != "" {
					return tag.JsonName, tag.Options
				}

				return strcase.ToLowerCamel(tag.FieldName), tag.Options
			}),
			livereload.New(livereload.Options{}),
		},
	})

	if err := g.Run(); err != nil {
		panic(err)
	}
}
