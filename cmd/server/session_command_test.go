package main

import (
	"reflect"
	"testing"

	"github.com/RandomChugokujin/CruxC2/internal/sessions"
)

func TestSessionChoicesUseHostnameAndKeepAgentIDsDistinct(t *testing.T) {
	agents := []*sessions.Session{
		{AgentID: "b-agent", Hostname: "srv-b", Platform: "Linux", Architecture: "amd64"},
		{AgentID: "a-agent", Hostname: "srv-a", Platform: "Windows", Architecture: "arm64"},
	}

	choices, idsByChoice := sessionChoices(agents)

	wantChoices := []string{
		"srv-a (Windows/arm64) [a-agent]",
		"srv-b (Linux/amd64) [b-agent]",
	}
	if !reflect.DeepEqual(choices, wantChoices) {
		t.Fatalf("sessionChoices() choices = %#v, want %#v", choices, wantChoices)
	}
	if got := idsByChoice["srv-a (Windows/arm64) [a-agent]"]; got != "a-agent" {
		t.Fatalf("selected choice resolves to %q, want %q", got, "a-agent")
	}
}

func TestSessionChoicesFallBackToAgentIDWhenHostnameIsEmpty(t *testing.T) {
	agents := []*sessions.Session{{AgentID: "agent-a", Platform: "Linux", Architecture: "amd64"}}
	choices, _ := sessionChoices(agents)
	want := []string{"agent-a (Linux/amd64)"}
	if !reflect.DeepEqual(choices, want) {
		t.Fatalf("sessionChoices() choices = %#v, want %#v", choices, want)
	}
}

func TestSessionPromptUsesHostname(t *testing.T) {
	got := sessionPrompt("srv-a", "a-agent")
	want := "CruxC2 [srv-a]> "
	if got != want {
		t.Fatalf("sessionPrompt() = %q, want %q", got, want)
	}
}

func TestSessionPromptFallsBackToAgentIDWhenHostnameIsEmpty(t *testing.T) {
	got := sessionPrompt("", "a-agent")
	want := "CruxC2 [a-agent]> "
	if got != want {
		t.Fatalf("sessionPrompt() = %q, want %q", got, want)
	}
}

func TestClearSelectedSessionAfterRemovalResetsPrompt(t *testing.T) {
	prompt := "CruxC2 [srv-a]> "
	setPrompt := func(value string) { prompt = value }
	selectedAgentID := "agent-a"

	clearSelectionOnRemoval(&selectedAgentID, "agent-a", setPrompt)

	if selectedAgentID != "" {
		t.Fatalf("selected agent ID = %q, want empty", selectedAgentID)
	}
	if want := "CruxC2> "; prompt != want {
		t.Fatalf("prompt = %q, want %q", prompt, want)
	}
}

func TestClearSelectedSessionAfterRemovalKeepsOtherSelection(t *testing.T) {
	prompt := "CruxC2 [srv-a]> "
	setPrompt := func(value string) { prompt = value }
	selectedAgentID := "agent-a"

	clearSelectionOnRemoval(&selectedAgentID, "agent-b", setPrompt)

	if selectedAgentID != "agent-a" {
		t.Fatalf("selected agent ID = %q, want %q", selectedAgentID, "agent-a")
	}
	if want := "CruxC2 [srv-a]> "; prompt != want {
		t.Fatalf("prompt = %q, want %q", prompt, want)
	}
}

func TestRemoveAgentSessionRemovesAndResetsPromptForSelectedAgent(t *testing.T) {
	manager := sessions.NewManager()
	manager.Add(&sessions.Session{AgentID: "agent-a"})
	prompt := "CruxC2 [srv-a]> "
	setPrompt := func(value string) { prompt = value }
	selectedAgentID := "agent-a"

	if !removeAgentSession(manager, &selectedAgentID, "agent-a", setPrompt) {
		t.Fatal("removeAgentSession() = false, want true")
	}
	if _, found := manager.Get("agent-a"); found {
		t.Fatal("removed agent remains in session manager")
	}
	if selectedAgentID != "" {
		t.Fatalf("selected agent ID = %q, want empty", selectedAgentID)
	}
	if want := "CruxC2> "; prompt != want {
		t.Fatalf("prompt = %q, want %q", prompt, want)
	}
}

func TestRemoveOtherAgentSessionPreservesSelectedPrompt(t *testing.T) {
	manager := sessions.NewManager()
	manager.Add(&sessions.Session{AgentID: "agent-a"})
	manager.Add(&sessions.Session{AgentID: "agent-b"})
	prompt := "CruxC2 [srv-a]> "
	setPrompt := func(value string) { prompt = value }
	selectedAgentID := "agent-a"

	if !removeAgentSession(manager, &selectedAgentID, "agent-b", setPrompt) {
		t.Fatal("removeAgentSession() = false, want true")
	}
	if selectedAgentID != "agent-a" {
		t.Fatalf("selected agent ID = %q, want %q", selectedAgentID, "agent-a")
	}
	if want := "CruxC2 [srv-a]> "; prompt != want {
		t.Fatalf("prompt = %q, want %q", prompt, want)
	}
}
