package tasks

type Task interface {
	Init() TaskState
	Run(TaskStateListener)
	Done()
}

type TaskStateListener func(state TaskState)

type TaskState struct {
	ID       int64
	TypeID   int64
	Name     string
	State    string
	Error    error
	Progress int
}
