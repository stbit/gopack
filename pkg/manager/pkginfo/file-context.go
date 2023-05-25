package pkginfo

import (
	"go/ast"
	"go/token"
)

type FileContext struct {
	sourcePath string
	distPath   string
	ModuleName string
	Fset       *token.FileSet
	File       *ast.File
}

func (f *FileContext) GetSourcePath() string {
	return f.sourcePath
}

func (f *FileContext) GetDistPath() string {
	return f.distPath
}
