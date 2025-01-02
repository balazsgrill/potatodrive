package tasks

type TaskExecutor struct {
	tasks    chan Task
	cancel   chan bool
	listener TaskStateListener
}

func NewTaskExecutor(listener TaskStateListener) *TaskExecutor {
	return &TaskExecutor{
		tasks:    make(chan Task, 1000),
		cancel:   make(chan bool),
		listener: listener,
	}
}

// AddTask attempts to insert a task into the executor.
// May be blocking if task queue is full
func (e *TaskExecutor) AddTask(task Task) {
	// TODO check if executor is still running
	state := task.Init()
	e.tasks <- task
	e.listener(state)
}

func (e *TaskExecutor) Run() {
	defer close(e.tasks)
	for {
		select {
		case task := <-e.tasks:
			task.Run(e.listener)
			task.Done()
		case <-e.cancel:
			// TODO stop task loop gracefully
			return
		}
	}
}
