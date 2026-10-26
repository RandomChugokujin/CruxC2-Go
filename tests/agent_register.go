package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/RandomChugokujin/CruxC2/internal/protocol"
)

func main(){
	agent_id := protocol.NewID()
	payload := protocol.MsgAgentRegister {
		AgentVersion: "0.1.0",
		Platform: "Linux",
		Architecture: "x86-64",
		Signature: "",
	}
	env, err := protocol.NewEnvelope(agent_id, payload)
	if err != nil {
		panic(err)
	}

	env_json, err := json.Marshal(env)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Post("http://127.0.0.1:8888/agent", "application/json", bytes.NewBuffer(env_json))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf(string(resp_body))

	// 2nd registration attempt with the same ID.
	resp2, err := client.Post("http://127.0.0.1:8888/agent", "application/json", bytes.NewBuffer(env_json))
	if err != nil {
		panic(err)
	}
	defer resp2.Body.Close()
	resp2_body, err := io.ReadAll(resp2.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf(string(resp2_body))
}
