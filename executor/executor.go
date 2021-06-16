package executor

// Executor is simply a function that takes input strings from a channel until closed,
// does something with them, and sends one or more (possibly nil) errors to the returned channel
// when done. The channel is closed at the end.
//
// The opts are options and are defined on a case by case basis, as needed. May be
// be empty, again, as defined individually.
type Executor func(inputs <-chan string, opts ...string) <-chan error
