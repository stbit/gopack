package manager

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/stbit/gopack/pkg/fsnotify"
	"github.com/stbit/gopack/pkg/manager/logger"
	"github.com/stbit/gopack/pkg/manager/pkginfo"
	"golang.org/x/exp/slices"
)

func (m *Manager) GetSourceFiles() []*pkginfo.FileInfo {
	return m.sourceFiles
}

func (m *Manager) loadSourceFiles() error {
	files := []string{}

	err := filepath.Walk(m.rootPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.Contains(path, "dist"+string(os.PathSeparator)) {
			files = append(files, path)
		}

		return nil
	})
	if err != nil {
		return err
	}

	for _, v := range files {
		exist := slices.ContainsFunc(m.sourceFiles, func(f *pkginfo.FileInfo) bool {
			return f.GetSourcePath() == v
		})

		if !exist {
			f := pkginfo.NewFileInfo(m.ModuleName, m.rootPath, v)

			if f.Error != nil {
				return f.Error
			} else {
				m.sourceFiles = append(m.sourceFiles, f)
			}
		}
	}

	return nil
}

func (m *Manager) Watch() {
	fmt.Println(logger.Magenta("start watching..."))
	fsnotify.New(m.rootPath, func(deletedFiles []string) {
		for _, v := range deletedFiles {
			i := slices.IndexFunc(m.sourceFiles, func(fi *pkginfo.FileInfo) bool {
				return fi.GetSourcePath() == v
			})

			if i != -1 {
				s := m.sourceFiles[i]
				if err := os.Remove(s.GetDistPath()); err != nil {
					logger.Error(err)
				}

				m.sourceFiles = slices.Delete(m.sourceFiles, i, i+1)
			}
		}

		if err := m.parse(); err != nil {
			logger.Error(err)
		}
	})
}
