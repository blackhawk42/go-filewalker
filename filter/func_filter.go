package filter

// FilterFunc is a function used for string testing. It takes a string and
// returns a boolean saying if the string passed the test.
type FilterFunc func(str string) bool

// FunctionFilter is a Filter that performs a generic test on a stream of strings,
// based on a function.
type FunctionFilter struct {
	filterFunc FilterFunc
	done       chan struct{}
}

// NewFunctionFilter creates a new FunctionFilter based on a given FilterFunc.
func NewFunctionFilter(filterFunc FilterFunc) *FunctionFilter {
	return &FunctionFilter{
		filterFunc: filterFunc,
		done:       make(chan struct{}),
	}
}

// Start starts a filtering worker using the function given at creation.
//
// May be called multiple times for creating multiple workers.
func (filter *FunctionFilter) Start(inputs <-chan string, outputs chan<- string) {
	go func() {
		for s := range inputs {
			if filter.filterFunc(s) {
				outputs <- s
			}
		}

		filter.done <- struct{}{}
	}()
}

// Wait hangs until a worker started by Start returns.
//
// If Start was called multiple times, Wait should be called an equal number
// of times.
func (filter *FunctionFilter) Wait() {
	<-filter.done
}
