package pkginfo

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/hashicorp/go-multierror"
)

type FileInfo struct {
	sourcePath string
	distPath   string
	ModuleName string
	Fset       *token.FileSet
	File       *ast.File
	Error      error
}

func NewFileInfo(moduleName string, rootPath string, path string) *FileInfo {
	f := &FileInfo{
		sourcePath: path,
		distPath:   strings.Replace(path, rootPath, rootPath+string(os.PathSeparator)+"dist", 1),
		ModuleName: moduleName,
		Fset:       token.NewFileSet(),
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
	f.Error = multierror.Append(f.Error, err)
}
