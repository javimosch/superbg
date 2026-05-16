# superbg - Super Background Process Manager

## Overview
A Go CLI tool to run, track, and manage background processes on Linux. Acts as a lightweight job manager with persistent state.

## Commands

| Command | Description |
|---------|-------------|
| `superbg run <cmd> [args...]` | Run a command in the background, detach from terminal, capture output |
| `superbg list` | List all tracked processes with status |
| `superbg stop <id\|pid>` | Gracefully stop a process (SIGTERM) |
| `superbg kill <id\|pid>` | Force kill a process (SIGKILL) |
| `superbg logs <id\|pid>` | Tail/follow logs of a process |
| `superbg status <id\|pid>` | Show detailed status of a process |
| `superbg attach <id\|pid>` | Re-attach to a process's stdout/stderr |

## State Persistence

- State directory: `~/.superbg/`
- State file: `~/.superbg/state.json`
- Log files: `~/.superbg/logs/<id>.log`
- PID files: `~/.superbg/pids/<id>.pid`

### state.json format
```json
{
  "jobs": [
    {
      "id": 1,
      "name": "my-command",
      "command": ["python", "server.py"],
      "pid": 12345,
      "status": "running",
      "started_at": "2026-05-16T10:00:00Z",
      "stopped_at": "",
      "exit_code": 0
    }
  ],
  "next_id": 2
}
```

## Project Structure

```
/root/superbg/
├── main.go                  # Entry point, CLI parsing
├── cmd/
│   ├── run.go               # Run command implementation
│   ├── list.go              # List command implementation
│   ├── stop.go              # Stop command implementation
│   ├── kill.go              # Kill command implementation
│   ├── logs.go              # Logs command implementation
│   ├── status.go            # Status command implementation
│   └── attach.go            # Attach command implementation
├── state/
│   └── state.go             # State management (read/write JSON, path mgmt)
├── process/
│   └── process.go           # Process lifecycle (spawn, signal, monitor)
└── docs/plan/
    └── superbg-plan.md      # This plan file
```

## Key Design Decisions

1. **No external dependencies** — Use only the Go standard library.
2. **PID-based identification** — Track by PID, reference by auto-incrementing ID.
3. **Setsid for detach** — Use `syscall.Setsid()` to create a new session, detaching from the terminal.
4. **Output redirection** — Stdout/stderr go to `~/.superbg/logs/<id>.log`.
5. **Process reaping** — Background goroutine monitors child processes and updates state on exit.
6. **Graceful shutdown** — Handle SIGINT/SIGTERM to clean up state before exiting.

## Implementation Order

1. Project scaffolding — `main.go` with CLI arg parsing, state module
2. `state/state.go` — Read/write JSON state, directory initialization
3. `process/process.go` — Spawn detached process, signal handling, reaping
4. `cmd/run.go` — Wire up run command
5. `cmd/list.go` — Wire up list command
6. `cmd/stop.go`, `cmd/kill.go` — Signal-based process control
7. `cmd/logs.go`, `cmd/status.go`, `cmd/attach.go` — Monitoring commands
