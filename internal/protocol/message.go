package protocol

type MsgAgentRegister struct {
	AgentVersion 	string 		`json:"agent_version"`
	Platform		string 		`json:"platform"`
	Architecture	string 		`json:"architecture"`
	Signature 		string 		`json:"signature"`
}

type MsgAgentRegisterAck struct {
	AgentRegistered bool 	`json:"agent_registered"`
}

type MsgTaskAssign struct {
	TaskID		string 		`json:"task_id"`
	TaskType 	string 		`json:"task_type"`
	Params 		[]string	`json:"params"`
}

type MsgTaskResult struct {
	TaskID		string 		`json:"task_id"`
	TaskType 	string 		`json:"task_type"`
	TaskOutput 	string		`json:"task_output"`
}
