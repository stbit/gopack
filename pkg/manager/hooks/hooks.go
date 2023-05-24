package hooks

import (
	"log"

	"github.com/stbit/gopack/pkg/manager/pkginfo"
)

type (
	HookName      string
	ParseHookName string
)

var (
	HOOK_BEFORE_COMPILE HookName = "before:compile"
	HOOK_AFTER_COMPILE  HookName = "after:compile"

	HOOK_PARSE_FILE ParseHookName = "parse:file"
)

type Hooks interface {
	AddHook(pluginName string, name HookName, callback func() error)
	AddHookParseFile(pluginName string, name ParseHookName, callback func(f *pkginfo.FileInfo) error)
}

type ManagerHooks struct {
	hooks      []hook
	parseHooks []parseHook
}

type hook struct {
	PluginName string
	Name       HookName
	Callback   func() error
}

type parseHook struct {
	PluginName string
	Name       ParseHookName
	Callback   func(*pkginfo.FileInfo) error
}

func NewManager() *ManagerHooks {
	return &ManagerHooks{
		hooks:      make([]hook, 0),
		parseHooks: make([]parseHook, 0),
	}
}

func (m *ManagerHooks) AddHook(pluginName string, name HookName, callback func() error) {
	m.hooks = append(m.hooks, hook{PluginName: pluginName, Name: name, Callback: callback})
}

func (m *ManagerHooks) AddHookParseFile(pluginName string, name ParseHookName, callback func(*pkginfo.FileInfo) error) {
	m.parseHooks = append(m.parseHooks, parseHook{PluginName: pluginName, Name: name, Callback: callback})
}

func (m *ManagerHooks) EmitHook(name HookName) error {
	for _, v := range m.hooks {
		if v.Name == name {
			if err := v.Callback(); err != nil {
				log.Fatal(err)
			}
		}
	}

	return nil
}

func (m *ManagerHooks) EmitParseHook(name ParseHookName, f *pkginfo.FileInfo) error {
	for _, v := range m.parseHooks {
		if v.Name == name {
			if err := v.Callback(f); err != nil {
				log.Fatal(err)
			}
		}
	}

	return nil
}
