package main

import (
	"fmt"
	"sort"

	"github.com/RandomChugokujin/CruxC2/internal/sessions"
)

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
		choice := fmt.Sprintf("%s (%s/%s)", agent.AgentID, agent.Platform, agent.Architecture)
		choices = append(choices, choice)
		idsByChoice[choice] = agent.AgentID
	}
	return choices, idsByChoice
}
