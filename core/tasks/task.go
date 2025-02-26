package tasks

type Task interface {
	Init(ID uint64) TaskState
	Run(TaskStateListener)
	Done()
}

type TaskStateListener func(state TaskState)

type TaskState struct {
	ID       uint64
	TypeID   int64
	Name     string
	State    string
	Error    error
	Progress int
}

type BaseTask struct {
	ID     uint64
	TypeID int64
}

func (bt *BaseTask) Init(ID uint64) TaskState {
	bt.ID = ID
	return TaskState{ID: ID}
}

type TaskContext interface {
	SetName(name string)
	SetState(state string)
	SetProgress(progress int)
	TaskConsumer
}

type AtomicTask func(TaskContext) error

type WrappedTask struct {
	AtomicTask
	TaskConsumer
	State    TaskState
	listener TaskStateListener
}

func (wt *WrappedTask) Init(ID uint64) TaskState {
	wt.State.ID = ID
	return wt.State
}

func (wt *WrappedTask) SetName(name string) {
	wt.State.Name = name
	wt.listener(wt.State)
}

func (wt *WrappedTask) SetState(state string) {
	wt.State.State = state
	wt.listener(wt.State)
}

func (wt *WrappedTask) SetProgress(progress int) {
	wt.State.Progress = progress
	wt.listener(wt.State)
}

func (wt *WrappedTask) Run(listener TaskStateListener) {
	wt.listener = listener
	err := wt.AtomicTask(wt)
	if err != nil {
		wt.State.Error = err
		listener(wt.State)
	}
}

func (wt *WrappedTask) Done() {

}

func AddAtomicTask(taskContext TaskContext, task AtomicTask, typeID int64, name string) {
	taskContext.AddTask(WrapTask(taskContext, task, typeID, name))
}

func WrapTask(executor TaskConsumer, task AtomicTask, typeID int64, name string) Task {
	return &WrappedTask{
		TaskConsumer: executor,
		AtomicTask:   task,
		State: TaskState{
			TypeID: typeID,
			Name:   name,
		},
	}
}
