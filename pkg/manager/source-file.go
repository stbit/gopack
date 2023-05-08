package manager

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type SourceFile struct {
	sourcePath string
	distPath   string
}

func newSourceFile(rootPath string, sourcePath string) *SourceFile {
	distPath := strings.Replace(sourcePath, rootPath, rootPath+"/dist", 1)

	return &SourceFile{
		sourcePath: sourcePath,
		distPath:   distPath,
	}
}

func loadSourceFiles(rootPath string) ([]*SourceFile, error) {
	r := []*SourceFile{}

	err := filepath.Walk(rootPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.Contains(path, "dist") {
			r = append(r, newSourceFile(rootPath, path))
		}

		return nil
	})

	return r, err
}

func (s *SourceFile) Parse() error {
	if err := os.MkdirAll(filepath.Dir(s.distPath), os.ModePerm); err != nil {
		return err
	}

	p := newParseFile(s.sourcePath)

	if err := p.parse(s.distPath); err != nil {
		return err
	}

	return nil
}
