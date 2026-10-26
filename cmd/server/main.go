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
	c2 := server.NewC2(sessions.NewManager())

	fmt.Println("Welcome to CruxC2. Type *help* to get a list of available commands. Happy Hacking!")

	// Create CLI App
	app := grumble.New(&grumble.Config{
		Name:        "CruxC2",
		Description: "C2 Server CLI",
		Prompt:      "CruxC2> ",
	})

	selectedAgentID := ""

	app.AddCommand(&grumble.Command{
		Name: "session",
		Help: "Select an active agent session",
		Run: func(c *grumble.Context) error {
			choices, idsByChoice := sessionChoices(c2.List())
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
			app.SetPrompt(fmt.Sprintf("CruxC2 [%s]> ", selectedAgentID))
			fmt.Printf("Selected session: %s\n", selectedAgentID)
			return nil
		},
	})

	app.AddCommand(&grumble.Command{
		Name: "agent_list",
		Help: "List all connect agents",
		Run: func(c *grumble.Context) error {
			sessions := c2.List()
			fmt.Print("Currently Active Agents:\n\n")
			for _, s := range sessions {
				s.PrintSession()
			}
			return nil
		},
	})
	app.AddCommand(&grumble.Command{
		Name: "agent_rm",
		Help: "Disconnect an agent",
		Run: func(c *grumble.Context) error {
			prompt := &survey.Select{
				Message: "Select an agent to delete:",
				Options: []string{"Agent 1", "Agent 2"},
			}
			var sel string
			if err := survey.AskOne(prompt, &sel); err != nil {
				return err
			}
			fmt.Printf("Removing listener: %s\n", sel)
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
			mux.HandleFunc("POST /agent", c2.Register)

			listener := &http.Server{
				Addr:    ":" + port,
				Handler: mux,
			}
			go listener.ListenAndServe()
			fmt.Printf("Starting listener at port %s...\n", port)
			return nil
		},
	})

	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
