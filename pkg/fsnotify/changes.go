package fsnotify

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/exp/slices"
)

func setupChangesFiles(rootPath string, w *fsnotify.Watcher, onChange func(d []string)) {
	go func() {
		ignoreFolders := []string{rootPath + string(os.PathSeparator) + "dist", rootPath + string(os.PathSeparator) + "tmp"}
		interval := 100 * time.Millisecond
		changed := false
		deletedFiles := []string{}

		for {
			select {
			case event, ok := <-w.Events:
				if !ok {
					return
				}

				if !slices.Contains(ignoreFolders, event.Name) {
					if event.Has(fsnotify.Create) {
						if s, err := os.Stat(event.Name); err == nil && s.IsDir() {
							w.Add(event.Name)
						}
					}

					if strings.HasSuffix(event.Name, ".go") {
						if event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) || event.Has(fsnotify.Write) {
							deletedFiles = append(deletedFiles, event.Name)
						}

						changed = true
					}
				}

			case err, ok := <-w.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)

			case <-time.After(interval):
				if changed {
					d := make([]string, len(deletedFiles))
					copy(d, deletedFiles)
					deletedFiles = []string{}
					changed = false
					onChange(d)
				}
			}
		}
	}()
}
