package manager

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/stbit/gopack/pkg/manager/parser"
)

type Manager struct {
	rootPath string
}

func NewManager(rootPath string) (*Manager, error) {
	return &Manager{rootPath: rootPath}, nil
}

func (m *Manager) loadSourceFiles() ([]string, error) {
	r := []string{}

	err := filepath.Walk(m.rootPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.Contains(path, "dist"+string(os.PathSeparator)) {
			r = append(r, path)
		}

		return nil
	})

	return r, err
}

func (m *Manager) parse() error {
	distPath := m.rootPath + string(os.PathSeparator) + "dist"
	if err := os.RemoveAll(distPath); err != nil {
		return err
	}

	l, err := m.loadSourceFiles()
	if err != nil {
		return err
	}

	parser.LoadPackages(l)

	return nil
}

func (m *Manager) Run() error {
	return m.parse()
}
