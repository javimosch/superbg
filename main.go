package main

import (
	"fmt"
	"os"

	"github.com/superbg/cli/cmd"
)

const usage = `superbg - Super Background Process Manager

Usage:
  superbg run [--watch] [--max-restarts=N] [--env-file FILE] <command> [args...]
                          Run a command in the background
  superbg list            List background processes
  superbg stop [--timeout N] <id|pid>
                          Stop a process (SIGTERM → SIGKILL after timeout)
  superbg kill <id|pid>   Kill a process (SIGKILL)
  superbg logs [--follow|-f] <id|pid>
                          Show process logs
  superbg status <id|pid> Show process status
  superbg attach <id|pid> Follow process logs in real-time
  superbg help            Show this help message

Flags:
  --watch         Auto-restart the process when it exits
  --max-restarts=N  Max restarts (default 10, with --watch)
  --env-file FILE Load environment variables from FILE
  --timeout N     Seconds to wait before SIGKILL (default 5)
`

func main() {
	if len(os.Args) < 2 {
		fmt.Print(usage)
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	var err error

	switch command {
	case "run":
		err = cmd.Run(args)
	case "list":
		err = cmd.List()
	case "stop":
		if len(args) == 0 {
			err = fmt.Errorf("missing arguments.\nUsage: superbg stop [--timeout N] <id|pid>")
		} else {
			err = cmd.Stop(args)
		}
	case "kill":
		if len(args) == 0 {
			err = fmt.Errorf("missing arguments.\nUsage: superbg kill <id|pid>")
		} else {
			err = cmd.Kill(args)
		}
	case "logs":
		err = handleLogs(args)
	case "status":
		if len(args) == 0 {
			err = fmt.Errorf("missing arguments.\nUsage: superbg status <id|pid>")
		} else {
			err = cmd.Status(args[0])
		}
	case "attach":
		if len(args) == 0 {
			err = fmt.Errorf("missing arguments.\nUsage: superbg attach <id|pid>")
		} else {
			err = cmd.Attach(args[0])
		}
	case "help", "--help", "-h":
		fmt.Print(usage)
		return
	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Print(usage)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func handleLogs(args []string) error {
	follow := false
	filtered := []string{}

	for _, a := range args {
		if a == "--follow" || a == "-f" {
			follow = true
		} else {
			filtered = append(filtered, a)
		}
	}

	if len(filtered) == 0 {
		return fmt.Errorf("missing id/pid.\nUsage: superbg logs [--follow|-f] <id|pid>")
	}

	return cmd.Logs(filtered[0], follow)
}
