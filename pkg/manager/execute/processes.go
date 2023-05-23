package execute

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

func StartProcesses(fl CommandsFlag) func() {
	ctxs := make([]*exec.Cmd, fl.Len())

	for i, v := range fl {
		cmd := startCmd(v)
		ctxs[i] = cmd
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
	}

	return func() {
		for _, c := range ctxs {
			if os.PathSeparator == '\\' {
				killWindow(c)
			} else {
				killOther(c)
			}
		}
	}
}

func startCmd(c string) *exec.Cmd {
	if os.PathSeparator == '\\' {
		splits := strings.Split(c, " ")
		return exec.Command(splits[0], splits[1:]...)
	} else {
		cmd := exec.Command("sh", "-c", c)
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

		return cmd
	}
}

func killWindow(cmd *exec.Cmd) error {
	kill := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(cmd.Process.Pid))
	kill.Stderr = os.Stderr
	kill.Stdout = os.Stdout

	return kill.Run()
}

func killOther(cmd *exec.Cmd) error {
	err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	_, _ = cmd.Process.Wait()
	return err
}
