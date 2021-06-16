package executor

import (
	"bufio"
	"container/heap"
	"fmt"
	"io"
	"os"
)

// StringHeap is a heap for strings, implementing heap.Interface.
//
// A lot basically copy-pasted from sort.StringSlice
type StringHeap []string

func (sheap StringHeap) Len() int           { return len(sheap) }
func (sheap StringHeap) Less(i, j int) bool { return sheap[i] < sheap[j] }
func (sheap StringHeap) Swap(i, j int)      { sheap[i], sheap[j] = sheap[j], sheap[i] }

func (sheap *StringHeap) Push(x interface{}) {
	*sheap = append(*sheap, x.(string))
}

func (sheap *StringHeap) Pop() interface{} {
	old := *sheap
	n := len(old)
	item := old[n-1]
	old[n-1] = ""
	*sheap = old[0 : n-1]

	return item
}

// ReportExecutor is an Executor that sorts received strings lexicographically and then writes them
// to the file in the first element of opts, if any.
//
// If there are no options, or the first option is the empty string "", os.Stdout
// is used. Otherwise, it's taken as a filename. The file will be created if
// non-existant, and overwritten otherwise.
func ReportExecutor(inputs <-chan string, opts ...string) <-chan error {
	done := make(chan error)

	filename := ""
	if len(opts) > 0 {
		filename = opts[0]
	}

	var outDevice io.Writer
	var err error
	if filename == "" {
		outDevice = os.Stdout
	} else {
		outDevice, err = os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			go func() {
				done <- fmt.Errorf("ReportExecutor: %v", err)
				close(done)
			}()

			return done
		}
	}

	go func() {
		defer close(done)

		sheap := make(StringHeap, 0)
		buff := bufio.NewWriter(outDevice)

		for in := range inputs {
			heap.Push(&sheap, in)
		}

		n := sheap.Len()
		for i := 0; i < n; i++ {
			buff.WriteString(heap.Pop(&sheap).(string))
			buff.WriteString("\n")
		}

		err := buff.Flush()
		if err != nil {
			done <- fmt.Errorf("ReportExecutor: %v", err)
			return
		}

		done <- nil
	}()

	return done
}
