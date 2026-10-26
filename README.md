# CruxC2-Go

> [!WARNING]
> This project is for educational or offensive security use in controlled, authorized environments only. Any misuse against unauthorized systems is illegal and strictly discouraged.

CruxC2-Go is a HTTP-based Command-and-Control (C2) framework for Offensive Cybersecurity.

## Build Instructions
Clone the repo and run `make`.
```
git clone https://github.com/RandomChugokujin/CruxC2-Go.git && cd CruxC2-Go
make
```

The `crux-server` binary can be found under the `bin/` directory:
```
$ ./bin/crux-server

  /$$$$$$                                           /$$$$$$
 /$$__  $$                                         /$$__  $$
| $$  \__/  /$$$$$$  /$$   /$$ /$$   /$$  /$$$$$$$|__/  \ $$
| $$       /$$__  $$| $$  | $$|  $$ /$$/ /$$_____/  /$$$$$$/
| $$      | $$ \__/| $$  | $$ \  $$$$/ | $$       /$$____/
| $$    $$| $$      | $$  | $$  >$$  $$ | $$      | $$
|  $$$$$$/| $$      |  $$$$$$/ /$$/\  $$|  $$$$$$$| $$$$$$$$
 \______/ |__/       \______/ |__/  \__/ \_______/|________/

Author: Brian Cao (RandomChugokujin)
Welcome to CruxC2. Type *help* to get a list of available commands. Happy Hacking!
CruxC2> help

C2 Server CLI

Commands:
  agent_list      List all connect agents
  agent_rm        Disconnect an agent
  clear           clear the screen
  exit            exit the shell
  help            use 'help [command]' for command help
  listener_start  Start HTTP listener on PORT_NUM
  session         Select an active agent session
```
