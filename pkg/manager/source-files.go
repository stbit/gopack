package manager

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/stbit/gopack/pkg/manager/pkginfo"
)

func (m *Manager) GetSourceFiles() []*pkginfo.FileInfo {
	return m.sourceFiles
}

func (m *Manager) loadSourceFiles() error {
	m.sourceFiles = []*pkginfo.FileInfo{}

	return filepath.Walk(m.rootPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.Contains(path, "dist"+string(os.PathSeparator)) {
			f := pkginfo.NewFileInfo(m.ModuleName, m.rootPath, path)

			if f.Error != nil {
				println(f.Error)
			} else {
				m.sourceFiles = append(m.sourceFiles, f)
			}
		}

		return nil
	})
}
