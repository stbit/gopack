package pkginfo

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
)

type FileInfo struct {
	*FileContext
	saved bool
}

func NewFileInfo(moduleName string, rootPath string, path string) *FileInfo {
	f := &FileInfo{
		FileContext: &FileContext{
			rootPath:   rootPath,
			sourcePath: path,
			distPath:   strings.Replace(path, rootPath, rootPath+string(os.PathSeparator)+"dist", 1),
			ModuleName: moduleName,
			Fset:       token.NewFileSet(),
			nodesLines: make(map[dst.Node]int),
		},
	}

	dec := decorator.NewDecoratorWithImports(f.Fset, "", nil)
	file, err := dec.ParseFile(path, nil, parser.ParseComments)
	if err != nil {
		f.AddError(err)
	}

	f.File = file
	f.initNodesNumberLine()
	return f
}

func (f *FileInfo) GetSourcePath() string {
	return f.sourcePath
}

func (f *FileInfo) GetDistPath() string {
	return f.distPath
}

func (f *FileInfo) IsSaved() bool {
	return f.saved
}

func (f *FileInfo) Save() error {
	// replace imports to dist
	for _, x := range f.File.Imports {
		if strings.HasPrefix(x.Path.Value, "\""+f.ModuleName) {
			x.Path.Value = strings.Replace(x.Path.Value, f.ModuleName, f.ModuleName+"/dist", 1)
		}
	}

	file := f.File
	distPath := f.GetDistPath()

	if err := os.MkdirAll(filepath.Dir(distPath), os.ModePerm); err != nil {
		panic(err)
	}

	of, err := os.OpenFile(distPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		panic(err)
	}

	defer of.Close()

	if err = decorator.Fprint(of, file); err != nil {
		return err
	}

	f.saved = true

	return nil
}

func (f *FileContext) AddImport(path string) {
	for i := 0; i < len(f.File.Decls); i++ {
		d := f.File.Decls[i]

		switch d.(type) {
		case *dst.FuncDecl:
			// No action
		case *dst.GenDecl:
			dd := d.(*dst.GenDecl)

			// IMPORT Declarations
			if dd.Tok == token.IMPORT {
				// Add the new import
				dd.Specs = append(dd.Specs, &dst.ImportSpec{
					Path: &dst.BasicLit{Kind: token.STRING, Value: strconv.Quote(path)},
				})
			}
		}
	}
}
