package pkginfo

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type FileInfo struct {
	*FileContext
	Error error
	saved bool
}

func NewFileInfo(moduleName string, rootPath string, path string) *FileInfo {
	f := &FileInfo{
		FileContext: &FileContext{
			sourcePath: path,
			distPath:   strings.Replace(path, rootPath, rootPath+string(os.PathSeparator)+"dist", 1),
			ModuleName: moduleName,
			Fset:       token.NewFileSet(),
		},
	}

	file, err := parser.ParseFile(f.Fset, path, nil, parser.ParseComments)
	if err != nil {
		f.AddError(err)
	}

	f.File = file
	return f
}

func (f *FileInfo) GetSourcePath() string {
	return f.sourcePath
}

func (f *FileInfo) GetDistPath() string {
	return f.distPath
}

func (f *FileInfo) AddError(err error) {
	f.Error = errors.Join(f.Error, err)
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

	ast.SortImports(f.Fset, f.File)
	f.File.Comments = []*ast.CommentGroup{}

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

	if err = printer.Fprint(of, f.Fset, file); err != nil {
		return err
	}

	f.saved = true

	return nil
}
