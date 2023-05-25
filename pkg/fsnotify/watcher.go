package fsnotify

import (
	"io/fs"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/stbit/gopack/pkg/manager/logger"
	"golang.org/x/exp/slices"
)

func New(rootPath string, onChange func(d []string)) {
	ignoreFolders := []string{"dist", "tmp"}
	w, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Fatal(err)
	}
	defer w.Close()

	setupChangesFiles(rootPath, w, onChange)

	if err := w.Add(rootPath); err != nil {
		logger.Fatal(err)
	}

	err = filepath.Walk(rootPath, func(path string, info fs.FileInfo, err error) error {
		filename := filepath.Base(path)
		isHidden := isHiddenFile(path)

		if info.IsDir() && (slices.Contains(ignoreFolders, filename) || isHidden) {
			return filepath.SkipDir
		}

		if info.IsDir() {
			return w.Add(path)
		}

		return nil
	})

	if err != nil {
		logger.Fatal(err)
	}

	<-make(chan struct{})
}

const dotCharacter = 46

func isHiddenFile(path string) bool {
	if filepath.Base(path)[0] == dotCharacter {
		return true
	}

	return false
}
