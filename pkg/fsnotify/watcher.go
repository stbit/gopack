package fsnotify

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/radovskyb/watcher"
)

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

	w.Ignore(ignoreDist)
	setupChangesFiles(w)

	if err := w.Add(rootPath); err != nil {
		log.Fatalln(err)
	}

	err := filepath.Walk(rootPath, func(path string, info fs.FileInfo, err error) error {
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

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
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
