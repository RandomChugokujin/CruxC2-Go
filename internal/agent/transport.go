package agent

import (
	"errors"
	"encoding/json"
	"github.com/RandomChugokujin/CruxC2/internal/protocol"
)

func NewRequest(msg any, agent_id string) ([]byte, error) {
	request_payload, err := json.Marshal(msg)
	if err != nil {
		return nil, errors.New("Payload mashalling failure")
	}
	request_env, err := protocol.NewEnvelope(agent_id, request_payload)
	if err != nil {
		return nil, errors.New("Envelope creation failure")
	}
	request_json, err := json.Marshal(request_env)
	if err != nil {
		return nil, errors.New("Envelope encoding failure")
	}
	return request_json, nil
}
