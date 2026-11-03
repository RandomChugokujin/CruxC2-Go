package protocol

import(
	"github.com/RandomChugokujin/CruxC2/internal/tasks"
)

type MsgAgentRegister struct {
	Hostname 		string 		`json:"hostname"`
	AgentVersion 	string 		`json:"agent_version"`
	Platform		string 		`json:"platform"`
	Architecture	string 		`json:"architecture"`
	Signature 		string 		`json:"signature"`
}

type MsgAgentRegisterAck struct {
	AgentRegistered bool 	`json:"agent_registered"`
}

type MsgTaskPoll struct {
	NumTasksExecuting 	int 	`json:"num_tasks_executing"`
}

type MsgTaskAssign struct {
	TaskID		string 				`json:"task_id"`
	TaskType 	tasks.TaskType 		`json:"task_type"`
	Params 		[]string			`json:"params"`
}

type MsgTaskResult struct {
	TaskID		string 		`json:"task_id"`
	TaskType 	string 		`json:"task_type"`
	TaskOutput 	string		`json:"task_output"`
}

type MsgTaskResultAck struct {
	TaskID		string 		`json:"task_id"`
}
