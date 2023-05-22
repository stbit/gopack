package fsnotify

import (
	"io/fs"
	"log"
	"path/filepath"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

func New(rootPath string, onChange func()) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	setupChangesFiles(rootPath, w, onChange)

	if err := w.Add(rootPath); err != nil {
		log.Fatalln(err)
	}

	err = filepath.Walk(rootPath, func(path string, info fs.FileInfo, err error) error {
		filename := filepath.Base(path)
		isHidden, err := isHiddenFile(path)
		if err != nil {
			return err
		}

		if info.IsDir() && (filename == "dist" || isHidden) {
			return filepath.SkipDir
		}

		if info.IsDir() {
			return w.Add(path)
		}

		return nil
	})

	if err != nil {
		log.Fatalln(err)
	}

	<-make(chan struct{})
}

func isHiddenFile(path string) (bool, error) {
	pointer, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return false, err
	}

	attributes, err := syscall.GetFileAttributes(pointer)
	if err != nil {
		return false, err
	}

	return attributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0, nil
}
