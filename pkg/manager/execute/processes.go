package execute

import (
	"log"
	"os"
	"os/exec"
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
		for _, v := range ctxs {
			if err := v.Process.Kill(); err != nil {
				log.Fatal(err)
			}
		}
	}
}
