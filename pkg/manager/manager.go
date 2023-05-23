package manager

import (
	"fmt"
	"go/ast"
	"go/printer"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/stbit/gopack/pkg/fsnotify"
	"github.com/stbit/gopack/pkg/manager/execute"
	"github.com/stbit/gopack/pkg/manager/logger"
	"github.com/stbit/gopack/pkg/manager/pkginfo"
	"github.com/stbit/gopack/pkg/plugins/syncerr"
	"golang.org/x/mod/modfile"
)

type Manager struct {
	mu             sync.Mutex
	rootPath       string
	ModuleName     string
	processManager *execute.ProcessManager
}

func New(rootPath string, fl execute.CommandsFlag) (*Manager, error) {
	modPath := rootPath + string(os.PathSeparator) + "go.mod"
	buf, err := ioutil.ReadFile(modPath)
	if err != nil {
		return nil, err
	}

	return &Manager{
		rootPath:       rootPath,
		ModuleName:     modfile.ModulePath(buf),
		processManager: execute.New(fl),
	}, nil
}

func (m *Manager) loadSourceFiles() ([]*pkginfo.FileInfo, error) {
	r := []*pkginfo.FileInfo{}

	err := filepath.Walk(m.rootPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.Contains(path, "dist"+string(os.PathSeparator)) {
			f := pkginfo.NewFileInfo(m.ModuleName, m.rootPath, path)

			if f.Error != nil {
				println(f.Error)
			} else {
				r = append(r, f)
			}
		}

		return nil
	})

	return r, err
}

func (m *Manager) clearDist() error {
	distPath := m.rootPath + string(os.PathSeparator) + "dist"
	if err := os.RemoveAll(distPath); err != nil {
		return err
	}

	return nil
}

func (m *Manager) parse() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	start := time.Now()

	if err := m.clearDist(); err != nil {
		return nil
	}

	l, err := m.loadSourceFiles()
	if err != nil {
		return err
	}

	for _, f := range l {
		syncerr.ParseFile(f)
	}

	for _, f := range l {
		if f.Error == nil {
			m.saveDistFile(f)
		}
	}

	log.Printf("compiled %s %s", logger.Success("successfully"), time.Since(start))

	m.processManager.Start()

	return nil
}

func (m *Manager) Run() error {
	if err := m.parse(); err != nil {
		return err
	}

	return nil
}

func (m *Manager) Watch() {
	fmt.Println("start watching...")
	fsnotify.New(m.rootPath, func() {
		if err := m.parse(); err != nil {
			log.Fatal(err)
		}
	})
}

func (p *Manager) saveDistFile(f *pkginfo.FileInfo) error {
	// replace imports to dist
	for _, x := range f.File.Imports {
		if strings.HasPrefix(x.Path.Value, "\""+p.ModuleName) {
			x.Path.Value = strings.Replace(x.Path.Value, p.ModuleName, p.ModuleName+"/dist", 1)
		}
	}

	ast.SortImports(f.Fset, f.File)

	file := f.File
	distPath := f.GetDistPath()

	if err := os.MkdirAll(filepath.Dir(distPath), os.ModePerm); err != nil {
		panic(err)
	}

	of, err := os.OpenFile(distPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o755)
	if err != nil {
		panic(err)
	}

	defer of.Close()

	if err = printer.Fprint(of, f.Fset, file); err != nil {
		return err
	}

	return nil
}
