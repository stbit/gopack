package manager

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/stbit/gopack/pkg/manager/execute"
	"github.com/stbit/gopack/pkg/manager/hooks"
	"github.com/stbit/gopack/pkg/manager/logger"
	"github.com/stbit/gopack/pkg/manager/pkginfo"
	"github.com/stbit/gopack/plugins"
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

func New(rootPath string, watch bool, commands []execute.Command) (*Manager, error) {
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
		processManager: execute.New(commands),
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

	if err := os.Mkdir(m.distPath, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (m *Manager) parse() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	start := time.Now()

	if err := m.loadSourceFiles(); err != nil {
		return err
	}

	for _, f := range m.sourceFiles {
		if !f.IsSaved() {
			m.hooks.EmitParseHook(hooks.HOOK_PARSE_FILE, f.FileContext)
		}
	}

	for _, f := range m.sourceFiles {
		if f.Error == nil && !f.IsSaved() {
			f.Save()
		}
	}

	log.Printf("compiled %s %s", logger.Success("successfully"), time.Since(start))

	m.processManager.Start()

	return nil
}

func (m *Manager) Run() error {
	if err := m.clearDist(); err != nil {
		logger.Fatal(err)
	}

	mc := plugins.NewManagerContext(m.hooks)
	for _, v := range m.plugins {
		v.Register(mc)
	}

	if err := m.parse(); err != nil {
		logger.Error(err)
	}

	return nil
}
