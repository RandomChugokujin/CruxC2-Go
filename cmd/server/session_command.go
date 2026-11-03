package main

import (
	"fmt"
	"sort"

	"github.com/RandomChugokujin/CruxC2/internal/sessions"
)

const defaultPrompt = "CruxC2> "

// sessionChoices returns stable survey labels and a lookup from each label to
// its agent ID, so the UI text is never mistaken for the session identifier.
func sessionChoices(agents []*sessions.Session) ([]string, map[string]string) {
	ordered := append([]*sessions.Session(nil), agents...)
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].AgentID < ordered[j].AgentID
	})

	choices := make([]string, 0, len(ordered))
	idsByChoice := make(map[string]string, len(ordered))
	for _, agent := range ordered {
		displayName := agent.Hostname
		choice := ""
		if displayName == "" {
			displayName = agent.AgentID
			choice = fmt.Sprintf("%s (%s/%s)", displayName, agent.Platform, agent.Architecture)
		} else {
			choice = fmt.Sprintf("%s (%s/%s) [%s]", displayName, agent.Platform, agent.Architecture, agent.AgentID)
		}
		choices = append(choices, choice)
		idsByChoice[choice] = agent.AgentID
	}
	return choices, idsByChoice
}

func sessionPrompt(hostname, agentID string) string {
	if hostname == "" {
		hostname = agentID
	}
	return fmt.Sprintf("CruxC2 [%s]> ", hostname)
}

func clearSelectionOnRemoval(selectedAgentID *string, removedAgentID string, setPrompt func(string)) {
	if *selectedAgentID != removedAgentID {
		return
	}
	*selectedAgentID = ""
	setPrompt(defaultPrompt)
}

func removeAgentSession(manager sessions.SessionManager, selectedAgentID *string, agentID string, setPrompt func(string)) bool {
	if !manager.Remove(agentID) {
		return false
	}
	clearSelectionOnRemoval(selectedAgentID, agentID, setPrompt)
	return true
}
