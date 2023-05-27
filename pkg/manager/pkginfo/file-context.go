package pkginfo

import (
	"errors"
	"go/ast"
	"go/token"
)

type FileContext struct {
	sourcePath string
	distPath   string
	ModuleName string
	Error      error
	Fset       *token.FileSet
	File       *ast.File
}

func (f *FileContext) GetSourcePath() string {
	return f.sourcePath
}

func (f *FileContext) GetDistPath() string {
	return f.distPath
}

func (f *FileContext) AddError(err error) {
	f.Error = errors.Join(f.Error, err)
}
