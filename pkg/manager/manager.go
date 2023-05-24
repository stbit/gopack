package manager

import (
	"fmt"
	"go/ast"
	"go/printer"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/stbit/gopack/pkg/fsnotify"
	"github.com/stbit/gopack/pkg/manager/execute"
	"github.com/stbit/gopack/pkg/manager/hooks"
	"github.com/stbit/gopack/pkg/manager/logger"
	"github.com/stbit/gopack/pkg/manager/pkginfo"
	"github.com/stbit/gopack/pkg/plugins"
	"golang.org/x/mod/modfile"
)

type Manager struct {
	mu             sync.Mutex
	watch          bool
	rootPath       string
	distPath       string
	ModuleName     string
	processManager *execute.ProcessManager
	hooks          *hooks.ManagerHooks
	plugins        []plugins.PluginRegister
	sourceFiles    []*pkginfo.FileInfo
}

func New(rootPath string, watch bool, fl execute.CommandsFlag) (*Manager, error) {
	modPath := rootPath + string(os.PathSeparator) + "go.mod"
	buf, err := ioutil.ReadFile(modPath)
	if err != nil {
		return nil, err
	}

	return &Manager{
		rootPath:       rootPath,
		distPath:       rootPath + string(os.PathSeparator) + "dist",
		watch:          watch,
		ModuleName:     modfile.ModulePath(buf),
		processManager: execute.New(fl),
		hooks:          hooks.NewManager(),
		plugins:        make([]plugins.PluginRegister, 0),
		sourceFiles:    make([]*pkginfo.FileInfo, 0),
	}, nil
}

func (m *Manager) RegisterPlugin(p plugins.PluginRegister) {
	m.plugins = append(m.plugins, p)
}

func (m *Manager) clearDist() error {
	if err := os.RemoveAll(m.distPath); err != nil {
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

	if err := m.loadSourceFiles(); err != nil {
		return err
	}

	for _, f := range m.sourceFiles {
		m.hooks.EmitParseHook(hooks.HOOK_PARSE_FILE, f)
	}

	for _, f := range m.sourceFiles {
		if f.Error == nil {
			m.saveDistFile(f)
		}
	}

	log.Printf("compiled %s %s", logger.Success("successfully"), time.Since(start))

	m.processManager.Start()

	return nil
}

func (m *Manager) Run() error {
	mc := plugins.NewManagerContext(m.hooks)
	for _, v := range m.plugins {
		v.Register(mc)
	}

	if err := m.parse(); err != nil {
		return err
	}

	return nil
}

func (m *Manager) Watch() {
	fmt.Println(logger.Magenta("start watching..."))
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
