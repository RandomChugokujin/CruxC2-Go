package server

import (
	"net/http"
	"time"
	"encoding/json"
	"github.com/RandomChugokujin/CruxC2/internal/protocol"
	"github.com/RandomChugokujin/CruxC2/internal/sessions"
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

func (c *C2) Register(w http.ResponseWriter, r *http.Request) {
	// Decoding request
	var request protocol.Envelope
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Request Data Decode Failed (Point 1)", http.StatusBadRequest)
		return
	}

	if request.Type != "MsgAgentRegister" {
		http.Error(w, "Unexpected Message Type", http.StatusBadRequest)
		return
	}
	// Unmarshalling payload message
	var request_payload protocol.MsgAgentRegister
	if err := json.Unmarshal(request.Payload, &request_payload); err != nil {
		http.Error(w, "Request Data Decode Failed (Point 2)", http.StatusBadRequest)
		return
	}
	// TODO: Verify agent cryptographic signature before register
	// Register Agent
	register_timestamp := time.Now()
	new_session := &sessions.Session {
		AgentID: request.AgentID,
		IPAddr: r.RemoteAddr,
		Version: request_payload.AgentVersion,
		Platform: request_payload.Platform,
		Architecture: request_payload.Architecture,
		RegisterTime: register_timestamp,
		LastPoll: register_timestamp,
	}
	agent_registered := c.sessions.Add(new_session)
	// Prepare response message payload
	response_payload, err := json.Marshal(protocol.MsgAgentRegisterAck{
		AgentRegistered: agent_registered,
	})
	if err != nil {
		http.Error(w, "Payload mashalling failed", http.StatusInternalServerError)
		return
	}
	// Prepare resposne envelope
	response,err := protocol.NewEnvelope(request.AgentID, response_payload)
	if err != nil {
		http.Error(w, "Envelope marshalling failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "encode failed", http.StatusInternalServerError)
		return
	}
}
