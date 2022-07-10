package executors

type Executor struct {
	name  string
	ready chan bool
}

func New(name string) *Executor {
	return &Executor{
		name: name,
		ready: make(chan bool),
	}
}

func (e *Executor) Name() string {
	return e.name
}

func (e *Executor) Ready() <-chan bool {
	return e.ready
}

// ReadyUp signals that Executor is ready for the next execution cycle
func (e *Executor) ReadyUp() {
	go func() {
		e.ready <- true
	}()
}
