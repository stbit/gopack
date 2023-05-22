package execute

import (
	"context"
	"log"
	"os"
	"os/exec"
	"strings"
)

func StartProcesses(fl CommandsFlag) func() {
	ctxs := make([]context.CancelFunc, fl.Len())

	for i, v := range fl {
		ctx, cancel := context.WithCancel(context.Background())

		ctxs[i] = cancel
		splits := strings.Split(v, " ")
		cmd := exec.CommandContext(ctx, splits[0], splits[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}
	}

	return func() {
		for _, v := range ctxs {
			v()
		}
	}
}
