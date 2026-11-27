package agent

import (
	"sync"
	"github.com/RandomChugokujin/CruxC2/internal/tasks"
)

type TaskManager interface {
	Add(new_task *tasks.Task) bool
	Get(task_id string) (tt *TaskTracker, ok bool)
	QueueSize() int
	Remove(task_id string) (ok bool)
	WriteResult(task_id string, output *string, status int) bool
}

func NewTaskManager() TaskManager {
	return &TaskQueue{task_queue: make(map[string]TaskTracker)}
}

type TaskQueue struct {
	mu 			sync.RWMutex
	task_queue	map[string]TaskTracker
}

func (tq *TaskQueue) Add(new_task *tasks.Task) bool {
	tq.mu.Lock()
	defer tq.mu.Unlock()
	if _, found := tq.task_queue[new_task.TaskID]; found {
		return false
	}
	// Refuse to add TaskNone or TaskShutdown
	if new_task.TaskType == tasks.TaskNone ||
	new_task.TaskType == tasks.TaskShutdown ||
	new_task.TaskType == tasks.TaskRegister {
		return false
	}
	tq.task_queue[new_task.TaskID] = TaskTracker{
		Task: 			*new_task,
		Complete: 		false,
	}
	return true
}
func (tq *TaskQueue) Get(task_id string) (tt *TaskTracker, ok bool) {
	tq.mu.RLock()
	defer tq.mu.RUnlock()
	if tt, found := tq.task_queue[task_id]; found {
		return &tt, true
	}
	return nil, false
}

func (tq *TaskQueue) QueueSize() int {
	tq.mu.RLock()
	defer tq.mu.RUnlock()
	return len(tq.task_queue)
}

func (tq *TaskQueue) Remove(task_id string) bool {
	tq.mu.Lock()
	defer tq.mu.Unlock()
	if _, found := tq.task_queue[task_id]; found {
		delete(tq.task_queue, task_id)
		return true
	}
	return false
}

func (tq *TaskQueue) WriteResult(task_id string, output *string, status int) bool {
	tq.mu.Lock()
	defer tq.mu.Unlock()
	if tt, found := tq.task_queue[task_id]; found {
		tt.Complete = true
		tt.Result.Output = *output
		return true
	}
	return false
}

type TaskTracker struct {
	Task 			tasks.Task
	Result			TaskResult
	Complete 		bool
}

type TaskResult struct {
	Status			int
	Output			string
}
