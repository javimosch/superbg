package main

import (
	"fmt"
	"os"

	"github.com/superbg/cli/cmd"
)

const usage = `superbg - Super Background Process Manager

Usage:
  superbg run [--watch] [--max-restarts=N] [--env-file FILE]
              [--name NAME] [--cwd DIR] <command> [args...]
                          Run a command in the background
  superbg list [--json]   List background processes
  superbg stop [--timeout N] <id|pid>
                          Stop a process (SIGTERM -> SIGKILL after timeout)
  superbg kill <id|pid>   Kill a process (SIGKILL)
  superbg logs [--follow|-f] [--json] <id|pid>
                          Show process logs
  superbg status [--json] <id|pid>
                          Show process status
  superbg attach <id|pid> Follow process logs in real-time
  superbg clean           Remove all completed processes from tracking
  superbg rm <id|pid>     Remove a single process from tracking
  superbg completion <bash|zsh|fish>
                          Generate shell completion script
  superbg help            Show this help message

Flags:
  --watch         Auto-restart the process when it exits
  --max-restarts=N  Max restarts (default 10, with --watch)
  --env-file FILE Load environment variables from FILE
  --name NAME     Custom process name
  --cwd DIR       Working directory for the process
  --timeout N     Seconds to wait before SIGKILL (default 5)
  --json          Output in JSON format
  --follow, -f    Follow log output in real-time
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
		asJSON := hasFlag(args, "--json")
		err = cmd.List(asJSON)
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
		asJSON := hasFlag(args, "--json")
		target := lastArg(args)
		if target == "" {
			err = fmt.Errorf("missing id/pid.\nUsage: superbg status [--json] <id|pid>")
		} else {
			err = cmd.Status(target, asJSON)
		}
	case "attach":
		if len(args) == 0 {
			err = fmt.Errorf("missing arguments.\nUsage: superbg attach <id|pid>")
		} else {
			err = cmd.Attach(args[len(args)-1])
		}
	case "clean":
		err = cmd.Clean()
	case "rm":
		if len(args) == 0 {
			err = fmt.Errorf("missing arguments.\nUsage: superbg rm <id|pid>")
		} else {
			err = cmd.Rm(args[len(args)-1])
		}
	case "completion":
		if len(args) == 0 {
			err = fmt.Errorf("missing arguments.\nUsage: superbg completion <bash|zsh|fish>")
		} else {
			err = cmd.Completion(args[len(args)-1])
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

func hasFlag(args []string, flag string) bool {
	for _, a := range args {
		if a == flag {
			return true
		}
	}
	return false
}

func lastArg(args []string) string {
	for i := len(args) - 1; i >= 0; i-- {
		if !hasPrefixAny(args[i], "--", "-") {
			return args[i]
		}
	}
	return ""
}

func hasPrefixAny(s string, prefixes ...string) bool {
	for _, p := range prefixes {
		if len(s) >= len(p) && s[:len(p)] == p {
			return true
		}
	}
	return false
}

func handleLogs(args []string) error {
	follow := false
	asJSON := false
	target := ""

	for _, a := range args {
		switch a {
		case "--follow", "-f":
			follow = true
		case "--json":
			asJSON = true
		default:
			target = a
		}
	}

	if target == "" {
		return fmt.Errorf("missing id/pid.\nUsage: superbg logs [--follow|-f] [--json] <id|pid>")
	}

	return cmd.Logs(target, follow, asJSON)
}
