package executor

type Executor struct {
	name  string
	ready chan bool
}

func New(name string) *Executor {
	return &Executor{
		name:  name,
		ready: make(chan bool),
	}
}

func (e *Executor) Name() string {
	return e.name
}

// Ready returns a channel to check if the Executor is ready for execution
func (e *Executor) Ready() <-chan bool {
	return e.ready
}

// ReadyUp signals the Executor to become ready
func (e *Executor) ReadyUp() {
	go func() {
		e.ready <- true
	}()
}
