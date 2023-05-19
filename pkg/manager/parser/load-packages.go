package parser

import (
	"log"
	"path/filepath"
	"time"

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

func LoadPackages(l []string) error {
	start := time.Now()
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

	elapsed := time.Since(start)
	log.Printf("load packages %s", elapsed)

	start = time.Now()
	for _, p := range pkgs {
		p.Save()
	}

	elapsed = time.Since(start)
	log.Printf("save packages %s", elapsed)

	return nil
}
