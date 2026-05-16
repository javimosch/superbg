# superbg

**Super Background** — a zero-config CLI to run, track, and manage background processes on Linux.

```bash
superbg run python server.py    # Detach & run
superbg list                     # Show all processes
superbg stop 1                   # Graceful stop
superbg logs --follow 1          # Tail logs
```

## Why superbg?

| You could use… | But superbg… |
|---|---|
| `nohup cmd &` | Tracks PIDs, saves logs, lets you list/stop/status later |
| `tmux` / `screen` | Doesn't need a terminal multiplexer — one-shot `run` and done |
| `systemd --user` units | Zero config, no unit files, works from any directory |
| `pm2` | Works with **any** language, not just Node.js |
| `supervisor` / `s6` | No daemon, no config files, no learning curve |

## Commands

| Command | Description |
|---------|-------------|
| `superbg run <cmd> [args...]` | Run a command detached from the terminal |
| `superbg list` | List all tracked processes |
| `superbg stop <id\|pid>` | Send SIGTERM to a process |
| `superbg kill <id\|pid>` | Send SIGKILL to a process |
| `superbg status <id\|pid>` | Show detailed process info |
| `superbg logs [--follow\|-f] <id\|pid>` | View process logs |
| `superbg attach <id\|pid>` | Follow logs in real-time (`tail -f`) |

## How it works

1. `superbg run` spawns the command in a **new session** (`setsid`), fully detached from the terminal — no SIGHUP, no accidental Ctrl+C.
2. Stdout and stderr are captured to `~/.superbg/logs/<id>.log`.
3. The process PID and metadata are saved to `~/.superbg/state.json`.
4. All commands (`list`, `stop`, `kill`, `status`, `logs`) read from this state file, so processes survive reboots and terminal closures.

## Install

```bash
go install github.com/javimosch/superbg@latest
```

Or download a binary from the [releases page](https://github.com/javimosch/superbg/releases).

## Requirements

- Linux (uses `setsid` syscall)
- Go 1.22+ (to build)
