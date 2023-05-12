package parser

import (
	"go/ast"
	"go/printer"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/packages"
)

type SourcePackage struct {
	pkg *packages.Package
}

func (p *SourcePackage) Save() {
	for i := range p.pkg.GoFiles {
		file := p.pkg.Syntax[i]
		parseAstFile(p, file)
	}

	for i, f := range p.pkg.GoFiles {
		file := p.pkg.Syntax[i]
		distPath := strings.Replace(f, p.pkg.Module.Dir, p.pkg.Module.Dir+string(os.PathSeparator)+"dist", 1)
		p.saveDistFile(distPath, file)
	}
}

func (p *SourcePackage) saveDistFile(distPath string, file *ast.File) error {
	if err := os.MkdirAll(filepath.Dir(distPath), os.ModePerm); err != nil {
		panic(err)
	}

	of, err := os.OpenFile(distPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		panic(err)
	}

	defer of.Close()

	if err = printer.Fprint(of, p.pkg.Fset, file); err != nil {
		return err
	}

	return nil
}
