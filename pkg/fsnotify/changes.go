package fsnotify

import (
	"log"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

func setupChangesFiles(rootPath string, w *fsnotify.Watcher, onChange func()) {
	go func() {
		ignoreDist := rootPath + string(os.PathSeparator) + "dist"
		interval := 100 * time.Millisecond
		changed := false

		for {
			select {
			case event, ok := <-w.Events:
				if !ok {
					return
				}

				if event.Name != ignoreDist {
					if event.Has(fsnotify.Create) {
						s, err := os.Stat(event.Name)
						if err != nil {
							log.Fatal(err)
						} else if s.IsDir() {
							w.Add(event.Name)
						}
					}

					changed = true
				}

			case err, ok := <-w.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)

			case <-time.After(interval):
				if changed {
					changed = false
					onChange()
				}
			}
		}
	}()
}
