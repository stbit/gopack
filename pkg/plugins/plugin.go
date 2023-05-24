package plugins

import (
	"github.com/stbit/gopack/pkg/manager/hooks"
)

type ManagerContext struct {
	hooks.Hooks
}

func NewManagerContext(h hooks.Hooks) *ManagerContext {
	return &ManagerContext{Hooks: h}
}

type PluginRegister interface {
	Register(*ManagerContext) error
}
