package plugins

import (
	"github.com/stbit/gopack/pkg/manager/hooks"
)

type ManagerContext struct {
	hooks.Hooks
	watch bool
}

func NewManagerContext(h hooks.Hooks, w bool) *ManagerContext {
	return &ManagerContext{Hooks: h, watch: w}
}

func (m *ManagerContext) IsWatch() bool {
	return m.watch
}

type PluginRegister = func(*ManagerContext) error
