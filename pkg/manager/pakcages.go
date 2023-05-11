package manager

import (
	"fmt"
	"go/ast"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"golang.org/x/exp/slices"
	"golang.org/x/tools/go/packages"
	"honnef.co/go/tools/go/ast/astutil"
)

type SourcePackage struct {
	pkg *packages.Package
}

const loadMode = packages.NeedName |
	packages.NeedFiles |
	packages.NeedCompiledGoFiles |
	packages.NeedImports |
	packages.NeedDeps |
	packages.NeedTypes |
	packages.NeedSyntax |
	packages.NeedTypesInfo |
	packages.NeedModule

func loadPackages(l []*SourceFile) {
	pkgpaths := make([]string, 0)
	pkgs := make([]*SourcePackage, 0)

	for _, s := range l {
		path := filepath.Dir(s.sourcePath)

		if !slices.Contains(pkgpaths, path) {
			pkgpaths = append(pkgpaths, path)
		}
	}

	config := &packages.Config{Mode: loadMode}
	lprog, err := packages.Load(config, pkgpaths...)
	if err != nil {
		panic(err)
	}

	for _, p := range lprog {
		pkgs = append(pkgs, &SourcePackage{
			pkg: p,
		})
	}

	f := lprog[0]

	fmt.Println(lprog, f, f.EmbedFiles, f.GoFiles, f.Module.Path)

	pkgs[4].Save()
}

func (p *SourcePackage) Save() {
	spew.Dump(p.pkg.Types.Scope())
	for i := range p.pkg.GoFiles {
		file := p.pkg.Syntax[i]
		astutil.Apply(file, nil, func(c *astutil.Cursor) bool {
			n := c.Node()
			switch x := n.(type) {
			case *ast.FuncDecl:
				fmt.Println(x)
			}

			return true
		})
	}
}
