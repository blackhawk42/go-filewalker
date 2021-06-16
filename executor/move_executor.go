package executor

import (
	"fmt"
	"os"
	"path/filepath"
)

// MoveExecutor moves files from one location to another.
//
// inputs should be the filepaths of the files one wishes to move.
//
// The first opt will be created as a directory, if it doesn't already exist, and all
// input files moved to it with the same basename.
// If there are no opts, or the first opt is the empty string "", a new
// destination directory will be created in the current location called
// "filewalker-TIME", where TIME is the current Unix Time with nanosecond precision.
func MoveExecutor(inputs <-chan string, opts ...string) <-chan error {
	done := make(chan error)

	dir := ""
	if len(opts) > 0 {
		dir = opts[0]
	}

	if dir == "" {
		dir = getName()
	}

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		go func() {
			done <- fmt.Errorf("MoveExecutor: %v", err)
			close(done)
		}()

		return done
	}

	go func() {
		defer close(done)

		var newPath string
		for file := range inputs {
			// TODO: overwites files with same name
			newPath = filepath.Join(dir, filepath.Base(file))

			err = os.Rename(file, newPath)
			if err != nil {
				done <- fmt.Errorf("MoveExecutor: %v", err)
				return
			}
		}

		done <- nil
	}()

	return done
}
