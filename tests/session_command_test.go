package main

import (
	"reflect"
	"testing"

	"github.com/RandomChugokujin/CruxC2/internal/sessions"
)

func TestSessionChoicesAreSortedAndResolveToAgentIDs(t *testing.T) {
	agents := []*sessions.Session{
		{AgentID: "b-agent", Platform: "Linux", Architecture: "amd64"},
		{AgentID: "a-agent", Platform: "Windows", Architecture: "arm64"},
	}

	choices, idsByChoice := sessionChoices(agents)

	wantChoices := []string{
		"a-agent (Windows/arm64)",
		"b-agent (Linux/amd64)",
	}
	if !reflect.DeepEqual(choices, wantChoices) {
		t.Fatalf("sessionChoices() choices = %#v, want %#v", choices, wantChoices)
	}
	if got := idsByChoice["a-agent (Windows/arm64)"]; got != "a-agent" {
		t.Fatalf("selected choice resolves to %q, want %q", got, "a-agent")
	}
}

func TestSessionChoicesEmptyWhenNoAgents(t *testing.T) {
	choices, idsByChoice := sessionChoices(nil)
	if len(choices) != 0 || len(idsByChoice) != 0 {
		t.Fatalf("sessionChoices(nil) = (%#v, %#v), want empty results", choices, idsByChoice)
	}
}
