package fsnotify

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/radovskyb/watcher"
)

type Watcher struct{}

func New(rootPath string) {
	w := watcher.New()
	w.FilterOps(watcher.Create, watcher.Move, watcher.Rename, watcher.Write, watcher.Remove)
	ignoreDist := rootPath + string(os.PathSeparator)

	w.AddFilterHook(func(info os.FileInfo, fullPath string) error {
		if strings.Contains(fullPath, ignoreDist) {
			return watcher.ErrSkip
		}

		return nil
	})

	go func() {
		for {
			select {
			case event := <-w.Event:
				fmt.Println(event) // Print the event's info.
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.AddRecursive(rootPath); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("fa", rootPath)

	for path, f := range w.WatchedFiles() {
		fmt.Printf("%s: %s\n", path, f.Name())
	}

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}
