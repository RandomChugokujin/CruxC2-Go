package tasks

import (
	"fmt"
	"strings"

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
	var task_str strings.Builder
	fmt.Fprintf(&task_str, "TaskID:\t%s\n", t.TaskID)
	fmt.Fprintf(&task_str, "TaskType:\t%s\n", t.TaskType.String())
	fmt.Fprintln(&task_str, "Params:")
	for _, p  := range t.Params {
		fmt.Fprintf(&task_str, "\t%s\n", p)
	}
	return task_str.String()
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
