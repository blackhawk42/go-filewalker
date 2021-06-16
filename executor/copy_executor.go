package executor

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyExecutor copies files to a destination.
//
// inputs should be a stream of filepaths to the files one wishes to copy.
//
// The first opt will be created as a directory, if it doesn't already exist, and
// all files will be copied to it with the same basename (overwriting existing files).
// If there are no opts, or the first opt is the empty string "", a new
// destination directory will be created in the current location called
// "filewalker-TIME", where TIME is the current Unix Time with nanosecond precision.
func CopyExecutor(inputs <-chan string, opts ...string) <-chan error {
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
			done <- fmt.Errorf("CopyExecutor: %v", err)
			close(done)
		}()

		return done
	}

	go func() {
		defer close(done)

		// TODO: modify buffer size with opts
		buff := make([]byte, 32*1024)
		var destPath string
		var err error
		var fdst *os.File
		var fsrc *os.File

		for file := range inputs {
			// TODO: this overwrites files with the same name
			destPath = filepath.Join(dir, filepath.Base(file))

			// Wrapped in a closure to easily defer Close calls and not exhaust
			// file handles
			func(file string) {
				fdst, err = os.OpenFile(destPath, os.O_RDWR|os.O_CREATE, os.ModePerm)
				if err != nil {
					return
				}
				defer fdst.Close()

				fsrc, err = os.Open(file)
				if err != nil {
					return
				}
				defer fsrc.Close()

				_, err = io.CopyBuffer(fdst, fsrc, buff)
				if err != nil {
					return
				}
			}(file)

			// Break at first captured error
			if err != nil {
				break
			}
		}

		if err != nil {
			done <- fmt.Errorf("CopyExecutor: %v", err)
			return
		}

		done <- nil
	}()

	return done
}
