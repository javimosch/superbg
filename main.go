package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/superbg/cli/cmd"
)

const usage = `superbg - Super Background Process Manager

Usage:
  superbg run <command> [args...]    Run a command in the background
  superbg list                       List background processes
  superbg stop <id|pid>              Stop a process (SIGTERM)
  superbg kill <id|pid>              Kill a process (SIGKILL)
  superbg logs [--follow|-f] <id|pid>  Show process logs
  superbg status <id|pid>            Show process status
  superbg attach <id|pid>            Follow process logs in real-time
  superbg help                       Show this help message
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
		err = requireArgs(args, "stop <id|pid>")
		if err == nil {
			err = cmd.Stop(args[0])
		}
	case "kill":
		err = requireArgs(args, "kill <id|pid>")
		if err == nil {
			err = cmd.Kill(args[0])
		}
	case "logs":
		err = handleLogs(args)
	case "status":
		err = requireArgs(args, "status <id|pid>")
		if err == nil {
			err = cmd.Status(args[0])
		}
	case "attach":
		err = requireArgs(args, "attach <id|pid>")
		if err == nil {
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

func requireArgs(args []string, usageStr string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing arguments.\nUsage: superbg %s", usageStr)
	}
	return nil
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

	return cmd.Logs(strings.Join(filtered, ""), follow)
}
