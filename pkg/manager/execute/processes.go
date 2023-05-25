package execute

import (
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/stbit/gopack/pkg/manager/logger"
)

type ProcessManager struct {
	mu   sync.Mutex
	ctxs []*exec.Cmd
	fl   CommandsFlag
}

func New(fl CommandsFlag) *ProcessManager {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	m := &ProcessManager{
		fl: fl,
	}

	go func() {
		<-sigs
		m.stop()
		os.Exit(0)
	}()

	return m
}

func (m *ProcessManager) Start() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.stop()

	if m.fl.Len() > 0 {
		m.ctxs = make([]*exec.Cmd, m.fl.Len())

		for i, v := range m.fl {
			cmd := m.startCmd(v)
			m.ctxs[i] = cmd
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Start(); err != nil {
				logger.Error(err)
			}
		}
	}
}

func (m *ProcessManager) stop() {
	if m.ctxs == nil {
		return
	}

	for _, c := range m.ctxs {
		if os.PathSeparator == '\\' {
			m.killWindow(c)
		} else {
			m.killOther(c)
		}
	}

	m.ctxs = nil
}

func (m *ProcessManager) startCmd(c string) *exec.Cmd {
	if os.PathSeparator == '\\' {
		splits := strings.Split(c, " ")
		return exec.Command(splits[0], splits[1:]...)
	} else {
		cmd := exec.Command("sh", "-c", c)
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

		return cmd
	}
}

func (m *ProcessManager) killWindow(cmd *exec.Cmd) error {
	kill := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(cmd.Process.Pid))
	kill.Stderr = os.Stderr
	kill.Stdout = os.Stdout

	return kill.Run()
}

func (m *ProcessManager) killOther(cmd *exec.Cmd) error {
	err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	_, _ = cmd.Process.Wait()
	return err
}
