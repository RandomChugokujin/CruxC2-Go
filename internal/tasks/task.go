package tasks

type TaskType int

const (
	TaskRegister = iota
	TaskExec
	TaskFileUpload
	TaskFileDownload
	TaskShutdown
)

var taskMap = map[TaskType]string {
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
