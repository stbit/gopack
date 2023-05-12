package parser2

import (
	"fmt"
	"path/filepath"

	"golang.org/x/exp/slices"
	"golang.org/x/tools/go/packages"
)

const loadMode = packages.NeedName |
	packages.NeedFiles |
	// packages.NeedCompiledGoFiles |
	// packages.NeedImports |
	// packages.NeedDeps |
	packages.NeedTypes |
	packages.NeedSyntax |
	packages.NeedTypesInfo |
	packages.NeedModule

func LoadPackages(l []string) {
	pkgpaths := make([]string, 0)
	pkgs := make([]*SourcePackage, 0)

	for _, s := range l {
		path := filepath.Dir(s)

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

	pkgs[5].Save()
}
