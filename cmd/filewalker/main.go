package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {

	// Flag parsing
	// TODO: Add better usage message
	var (
		filterWorkersNum = flag.Int("workers", runtime.NumCPU(), "`number` of concurrent filtering workers; defaults to detected CPUs")
		filterMethod     = flag.String("filter", "glob", "`name` of the method to use for filename matching; possible values: "+strings.Join(AvaiableFilterMethods.GetMethods(), ", "))
		pattern          = flag.String("pattern", "*", "pattern `string` to use in matching")
		outFile          = flag.String("out", "", "output `filename`; may be a directory or a file, specifics depend on the chosen action")
		executorMethod   = flag.String("action", "report", "`name` of the action to use; possible values: "+strings.Join(AvaiableExecutorMethods.GetMethods(), ", "))
	)

	flag.Parse()

	// Stablish base directory
	var baseDir string
	var err error
	if flag.NArg() > 0 {
		baseDir = flag.Arg(0)

		info, err := os.Stat(baseDir)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				fmt.Fprintf(os.Stderr, "error: %s does not exist\n", baseDir)
				helpAndExit(1)
			} else {
				fmt.Fprintf(os.Stderr, "error while reading %s: %v\n", baseDir, err)
				helpAndExit(1)
			}
		}

		if !info.IsDir() {
			fmt.Fprintf(os.Stderr, "error: %s is not a directory\n", baseDir)
			helpAndExit(1)
		}
	} else {
		baseDir, err = os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not read current working directory: %v\n", err)
			helpAndExit(1)
		}
	}

	// Stablish filtering method
	filterCreator, ok := AvaiableFilterMethods[*filterMethod]
	if !ok {
		fmt.Fprintf(os.Stderr, "error: %s is not a valid filtering method\n", *filterMethod)
		helpAndExit(1)
	}

	mainFilter, err := filterCreator(*pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: while creating filter: %v\n", err)
		helpAndExit(1)
	}

	// Stablish executor method
	executorFunc, ok := AvaiableExecutorMethods[*executorMethod]
	if !ok {
		fmt.Fprintf(os.Stderr, "error: %s is not a valid action\n", *executorMethod)
		helpAndExit(1)
	}

	// Set up working chain
	inputs := make(chan string, *filterWorkersNum)
	filteredPaths := make(chan string, *filterWorkersNum)

	for i := 0; i < *filterWorkersNum; i++ {
		mainFilter.Start(inputs, filteredPaths)
	}

	execWait := executorFunc(filteredPaths, *outFile)

	// Start work and collect results
	go func() {
		filepath.WalkDir(baseDir, func(path string, d fs.DirEntry, err error) error {
			if d.Type().IsRegular() {
				inputs <- path
			}

			return nil
		})
		close(inputs)
	}()

	go func() {
		for i := 0; i < *filterWorkersNum; i++ {
			mainFilter.Wait()
		}
		close(filteredPaths)
	}()

	for err := range execWait {
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			helpAndExit(1)
		}
	}
}

func helpAndExit(status int) {
	flag.Usage()
	os.Exit(status)
}
