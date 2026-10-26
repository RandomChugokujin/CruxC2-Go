package protocol

import (
	"encoding/json"
	"reflect"
	"time"
)

type Envelope struct {
    ProtocolVersion string    		`json:"protocol_version"`
    AgentID       	string    		`json:"agent_id"`
    MessageID       string    		`json:"message_id"`
    Type            string    		`json:"type"`
    Timestamp       time.Time 		`json:"timestamp"`
    Payload         json.RawMessage `json:"payload"`
}

func NewEnvelope(agent_id string, msg any) (Envelope, error){
	payload, err := json.Marshal(msg)

	if err != nil {
        // Marshal failure is a protocol bug, not a transient error;
        // surface it rather than send a broken envelope.
        return Envelope{}, err
    }
    return Envelope{
        ProtocolVersion: "0.1.0",
        AgentID:         agent_id,
        MessageID:       NewID(),
        Type:            shortTypeName(reflect.TypeOf(msg)),
        Timestamp:       time.Now(),
        Payload:         payload,
    }, nil
}

func shortTypeName(t reflect.Type) string {
    if t.Name() != "" {
        return t.Name()
    }
    // Fall back for unnamed types (unlikely here, but safe).
    return t.String()
}
