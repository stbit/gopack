package manager

import (
	"go/ast"
	"go/printer"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/stbit/gopack/pkg/manager/pkginfo"
	"github.com/stbit/gopack/pkg/plugins/syncerr"
	"golang.org/x/mod/modfile"
)

type Manager struct {
	rootPath   string
	ModuleName string
}

func New(rootPath string) (*Manager, error) {
	modPath := rootPath + string(os.PathSeparator) + "go.mod"
	buf, err := ioutil.ReadFile(modPath)
	if err != nil {
		return nil, err
	}

	return &Manager{
		rootPath:   rootPath,
		ModuleName: modfile.ModulePath(buf),
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

func (m *Manager) parse() error {
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

	return nil
}

func (m *Manager) Run() error {
	distPath := m.rootPath + string(os.PathSeparator) + "dist"
	if err := os.RemoveAll(distPath); err != nil {
		return err
	}

	// fsnotify.New(m.rootPath)

	return m.parse()
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
