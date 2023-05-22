package execute

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func StartProcesses(fl CommandsFlag) func() {
	ctxs := make([]*exec.Cmd, fl.Len())

	for i, v := range fl {
		splits := strings.Split(v, " ")
		cmd := exec.Command(splits[0], splits[1:]...)
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

func killWindow(cmd *exec.Cmd) error {
	kill := exec.Command("TASKKILL", "/T", "/F", "/PID", strconv.Itoa(cmd.Process.Pid))
	kill.Stderr = os.Stderr
	kill.Stdout = os.Stdout

	return kill.Run()
}

func killOther(cmd *exec.Cmd) error {
	kill := exec.Command("kill", "-15", strconv.Itoa(cmd.Process.Pid))
	kill.Stderr = os.Stderr
	kill.Stdout = os.Stdout

	return kill.Run()
}
