package gopack

import (
	"errors"
	"fmt"

	"github.com/stbit/gopack/pkg/manager"
	"github.com/stbit/gopack/pkg/manager/execute"
	"github.com/stbit/gopack/pkg/plugins"
)

var (
	RootPathErr    = errors.New("config RootPath undefined")
	CommandNameErr = errors.New("command name undefined")
	CommandExecErr = errors.New("command exec undefined")
)

type Gopack struct {
	config *Config
}

type Config struct {
	RootPath string
	Watch    bool
	Commands []Command
	Plugins  []plugins.PluginRegister
}

type Command struct {
	Name string
	Exec string
}

func New(c *Config) *Gopack {
	return &Gopack{
		config: c,
	}
}

func (g *Gopack) Run() error {
	commands := make([]execute.Command, 0)

	if g.config.Commands != nil && len(g.config.Commands) > 0 {
		for _, c := range g.config.Commands {
			if c.Name == "" {
				return CommandNameErr
			}

			if c.Exec == "" {
				return fmt.Errorf("command exec undefined for %s", c.Name)
			}

			commands = append(commands, execute.Command{Name: c.Name, Exec: c.Exec})
		}
	}

	if g.config.RootPath == "" {
		return RootPathErr
	}

	m, err := manager.New(g.config.RootPath, g.config.Watch, commands)
	if err != nil {
		return err
	}

	if g.config.Plugins != nil {
		for _, p := range g.config.Plugins {
			m.RegisterPlugin(p)
		}
	}

	if err = m.Run(); err != nil {
		panic(err)
	}

	if g.config.Watch {
		m.Watch()
	}

	return nil
}
