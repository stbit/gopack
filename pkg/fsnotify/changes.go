package fsnotify

import (
	"fmt"
	"log"

	"github.com/radovskyb/watcher"
)

func setupChangesFiles(w *watcher.Watcher) {
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
}
