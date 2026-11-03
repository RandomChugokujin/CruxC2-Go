package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/RandomChugokujin/CruxC2/internal/server"
	"github.com/RandomChugokujin/CruxC2/internal/sessions"

	"github.com/AlecAivazis/survey/v2"
	"github.com/desertbit/grumble"
)

var banner = `
  /$$$$$$                                           /$$$$$$
 /$$__  $$                                         /$$__  $$
| $$  \__/  /$$$$$$  /$$   /$$ /$$   /$$  /$$$$$$$|__/  \ $$
| $$       /$$__  $$| $$  | $$|  $$ /$$/ /$$_____/  /$$$$$$/
| $$      | $$ \__/| $$  | $$ \  $$$$/ | $$       /$$____/
| $$    $$| $$      | $$  | $$  >$$  $$ | $$      | $$
|  $$$$$$/| $$      |  $$$$$$/ /$$/\  $$|  $$$$$$$| $$$$$$$$
 \______/ |__/       \______/ |__/  \__/ \_______/|________/

Author: Brian Cao (RandomChugokujin)
`

func main() {
	fmt.Print(banner)

	// Create C2 Manager
	sessionManager := sessions.NewManager()
	c2 := server.NewC2(sessionManager)
	listeners := newListenerRegistry()

	fmt.Println("Welcome to CruxC2. Type *help* to get a list of available commands. Happy Hacking!")

	// Create CLI App
	app := grumble.New(&grumble.Config{
		Name:        "CruxC2",
		Description: "C2 Server CLI",
		Prompt:      defaultPrompt,
	})

	selectedAgentID := ""

	app.AddCommand(&grumble.Command{
		Name: "session",
		Help: "Select an active agent session",
		Run: func(c *grumble.Context) error {
			agents := c2.List()
			choices, idsByChoice := sessionChoices(agents)
			if len(choices) == 0 {
				fmt.Println("No active agent sessions.")
				return nil
			}

			prompt := &survey.Select{
				Message: "Select an active agent session:",
				Options: choices,
			}
			var selectedChoice string
			if err := survey.AskOne(prompt, &selectedChoice); err != nil {
				return err
			}

			agentID, found := idsByChoice[selectedChoice]
			if !found {
				return fmt.Errorf("selected agent session %q is no longer available", selectedChoice)
			}

			selectedAgentID = agentID
			selectedHostname := ""
			for _, agent := range agents {
				if agent.AgentID == agentID {
					selectedHostname = agent.Hostname
					break
				}
			}
			app.SetPrompt(sessionPrompt(selectedHostname, selectedAgentID))
			fmt.Printf("Selected session: %s\n", selectedAgentID)
			return nil
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "agent_list",
		Help: "List all connect agents",
		Run: func(c *grumble.Context) error {
			sessions := c2.List()
			if len(sessions) == 0 {
				fmt.Print("No active agent sessions.\n")
				return nil
			}
			fmt.Print("Currently Active Agents:\n\n")
			for _, s := range sessions {
				s.PrintSession()
			}
			return nil
		},
	})
	app.AddCommand(&grumble.Command{
		Name: "agent_rm",
		Help: "Remove an agent session",
		Run: func(c *grumble.Context) error {
			choices, idsByChoice := sessionChoices(c2.List())
			if len(choices) == 0 {
				fmt.Println("No active agent sessions.")
				return nil
			}

			prompt := &survey.Select{
				Message: "Select an agent session to remove:",
				Options: choices,
			}
			var selectedChoice string
			if err := survey.AskOne(prompt, &selectedChoice); err != nil {
				return err
			}

			agentID, found := idsByChoice[selectedChoice]
			if !found {
				return fmt.Errorf("selected agent session %q is no longer available", selectedChoice)
			}
			if !removeAgentSession(sessionManager, &selectedAgentID, agentID, app.SetPrompt) {
				return fmt.Errorf("agent session %q is no longer registered", agentID)
			}

			fmt.Printf("Removed agent session: %s\n", agentID)
			return nil
		},
	})
	app.AddCommand(&grumble.Command{
		Name: "listener_start",
		Help: "Start HTTP listener on PORT_NUM",
		Args: func(a *grumble.Args) {
			a.String("port", "Port Number")
		},
		Run: func(c *grumble.Context) error {
			port := c.Args.String("port")
			// Create HTTP Listener
			mux := http.NewServeMux()
			mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("ok\n"))
			})
			// Agent Register
			mux.HandleFunc("POST /agent", c2.Register)
			// Agent Task Poll
			mux.HandleFunc("POST /poll", c2.TaskPoll)
			// Agent Task Result submission
			mux.HandleFunc("POST /result", c2.TaskSubmit)

			listenerAddr, err := listeners.start(":"+port, mux)
			if err != nil {
				return fmt.Errorf("could not start HTTP listener on port %s: %w", port, err)
			}
			fmt.Printf("Started HTTP listener on %s\n", listenerAddr)
			return nil
		},
	})
	app.AddCommand(&grumble.Command{
		Name: "listener_list",
		Help: "List running listeners",
		Run: func(c *grumble.Context) error {
			listener_str_arr := listeners.addresses()
			if len(listener_str_arr) == 0 {
				fmt.Println("No running listeners, use listener_start <PORT_NUM> to start one.")
				return nil
			}
			fmt.Println("Running Listners:")
			for index, addr := range listener_str_arr {
				fmt.Printf("Listener %d: %s\n", index, addr)
			}
			return nil
		},
	})


	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
