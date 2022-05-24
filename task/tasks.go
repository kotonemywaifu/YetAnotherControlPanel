package task

var tasks map[string]*Task

type Task struct {
	Name      string
	tick      int
	NeedTicks int
	Func      func()
}

func (t *Task) Tick() {
	t.tick++
	if t.tick >= t.NeedTicks {
		t.Func()
		t.tick = 0
	}
}

func AddTask(task *Task) {
	if tasks == nil {
		tasks = make(map[string]*Task)
	}
	tasks[task.Name] = task
}

func RemoveTask(name string) {
	if tasks == nil {
		return
	}
	delete(tasks, name)
}
