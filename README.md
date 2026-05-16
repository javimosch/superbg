# superbg

**Super Background** — a zero-config CLI to run, track, and manage background processes on Linux.

```bash
superbg run python server.py         # Detach & run
superbg run --watch ./server          # Auto-restart on crash
superbg run --env-file .env ./app     # Load env vars from file
superbg list                          # Show all processes
superbg stop --timeout 10 1           # Graceful stop with fallback kill
superbg logs --follow 1               # Tail logs
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
| `superbg run [flags] <cmd> [args...]` | Run a command detached from the terminal |
| `superbg list` | List all tracked processes |
| `superbg stop [--timeout N] <id\|pid>` | Send SIGTERM, then SIGKILL after N seconds |
| `superbg kill <id\|pid>` | Send SIGKILL immediately |
| `superbg status <id\|pid>` | Show detailed process info |
| `superbg logs [--follow\|-f] <id\|pid>` | View process logs |
| `superbg attach <id\|pid>` | Follow logs in real-time (`tail -f`) |

## Features

### 🔁 Auto-restart (`--watch`)

Run a process as a supervised service. If it crashes, superbg re-spawns it with exponential backoff (1s → 30s max).

```bash
superbg run --watch --max-restarts=5 ./server
```

- `--max-restarts=N` — limit restarts before giving up (default 10)
- `stop` signals the monitor process, which forwards SIGTERM to the child for graceful shutdown
- `kill` terminates both monitor and child
- `list` and `status` show restart counts

### ⏱ Graceful stop with timeout

```bash
superbg stop --timeout 15 1
```

Sends SIGTERM and waits up to 15 seconds. If the process is still alive, sends SIGKILL. Default timeout is 5 seconds.

### 📄 Environment file support

```bash
superbg run --env-file .env ./app
```

Load `KEY=VALUE` pairs from a file (supports `#` comments and `export` prefix) and inject them into the child's environment.

## How it works

1. `superbg run` spawns the command in a **new session** (`setsid`), fully detached from the terminal — no SIGHUP, no accidental Ctrl+C.
2. Stdout and stderr are captured to `~/.superbg/logs/<id>.log`.
3. The process PID and metadata are saved to `~/.superbg/state.json`.
4. When `--watch` is used, superbg stays alive as a monitor, re-spawning the child on exit.
5. All commands (`list`, `stop`, `kill`, `status`, `logs`) read from the state file, so processes survive reboots and terminal closures.

## Install

```bash
go install github.com/javimosch/superbg@latest
```

Or download a binary from the [releases page](https://github.com/javimosch/superbg/releases).

## Requirements

- Linux (uses `setsid` syscall)
- Go 1.22+ (to build)
