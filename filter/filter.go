package filter

// Filter is something that takes string inputs, evaluates them based on some rule,
// and outputs strings that pass.
type Filter interface {
	// Start starts a concurrent worker that takes a channel of inputs, evaluates them,
	// and sends passing inputs to a channel of outputs.
	//
	// May be used multiple times to start multiple workers, if the implementation
	// allows it. This should be clearly stated by the documentation.
	Start(inputs <-chan string, outputs chan<- string)

	// Wait should hang until one worker started by Start returns.
	//
	// If the implementation allows Start to be called multiple times, Wait should
	// be called an equal number of times.
	Wait()
}
