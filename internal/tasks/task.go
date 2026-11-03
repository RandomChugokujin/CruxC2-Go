package tasks

import (
	"fmt"

	"github.com/RandomChugokujin/CruxC2/internal/id"
)

type TaskType int

const (
	TaskNone = iota
	TaskRegister
	TaskExec
	TaskFileUpload
	TaskFileDownload
	TaskShutdown
)

var taskMap = map[TaskType]string {
	TaskNone:			"none",
	TaskRegister: 		"register",
	TaskExec: 			"exec",
	TaskFileUpload: 	"file_upload",
	TaskFileDownload: 	"file_download",
	TaskShutdown:		"shutdown",
}

func (ts TaskType) String() string {
	return taskMap[ts]
}

type Task struct {
	TaskID 		string
	TaskType 	TaskType
	Params 		[]string
}

func (t Task) String() string {
	task_str := ""
	task_str += fmt.Sprintf("TaskID:\t%s\n", t.TaskID)
	task_str += fmt.Sprintf("TaskType:\t%s\n", t.TaskType.String())
	task_str += fmt.Sprintln("Params:")
	for _, p  := range t.Params {
		task_str += fmt.Sprintf("\t%s\n", p)
	}
	return task_str
}

func NewTask(tt TaskType, params *[]string) Task {
	if params == nil {
		return Task{
			TaskID: id.NewID(),
			TaskType: tt,
			Params: []string{},
		}
	}
	return Task{
		TaskID: id.NewID(),
		TaskType: tt,
		Params: *params,
	}
}
