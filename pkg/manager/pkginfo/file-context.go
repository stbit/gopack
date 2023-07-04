package pkginfo

import (
	"errors"
	"go/token"
	"os"
	"strings"

	"github.com/dave/dst"
)

type FileContext struct {
	rootPath   string
	sourcePath string
	distPath   string
	ModuleName string
	Error      error
	Fset       *token.FileSet
	File       *dst.File
	nodesLines map[dst.Node]int
}

func (f *FileContext) GetSourcePath() string {
	return f.sourcePath
}

func (f *FileContext) GetRelativeSourcePath() string {
	return strings.Replace(f.sourcePath, f.rootPath+string(os.PathSeparator), "", 1)
}

func (f *FileContext) GetDistPath() string {
	return f.distPath
}

func (f *FileContext) AddError(err error) {
	f.Error = errors.Join(f.Error, err)
}
