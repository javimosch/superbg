# superbg

Super Background -- a zero-config CLI to run, track, and manage background processes on Linux.

```bash
superbg run python server.py              # Detach and run
superbg run --watch --max-restarts=5 ./server  # Auto-restart on crash
superbg run --env-file .env --name api --cwd /app ./entrypoint
superbg list --json                       # Machine-readable listing
superbg stop --timeout 10 1               # Graceful stop with fallback kill
superbg logs -f 1                         # Tail logs
superbg clean                             # Remove completed jobs
superbg completion bash > /etc/bash_completion.d/superbg
```

## Why superbg?

| Instead of | superbg gives you |
|---|---|
| nohup cmd & | PID tracking, log capture, list/stop/status later |
| tmux / screen | No terminal multiplexer needed -- one-shot run |
| systemd --user units | Zero config, no unit files, works from any directory |
| pm2 | Any language, not just Node.js |
| supervisor / s6 | No daemon, no config files, no learning curve |

## Commands

| Command | Description |
|---|---|
| superbg run [flags] <cmd> [args...] | Run a command detached from the terminal |
| superbg list [--json] | List tracked processes |
| superbg stop [--timeout N] <id\|pid> | SIGTERM, then SIGKILL after N seconds (default 5) |
| superbg kill <id\|pid> | SIGKILL immediately |
| superbg status [--json] <id\|pid> | Show detailed process info |
| superbg logs [--follow\|-f] [--json] <id\|pid> | View process logs |
| superbg attach <id\|pid> | Follow logs in real-time |
| superbg clean | Remove all completed processes from tracking |
| superbg rm <id\|pid> | Remove a single process from tracking |
| superbg completion <bash\|zsh\|fish> | Generate shell completion script |
| superbg help | Show this help message |

## Run flags

| Flag | Description |
|---|---|
| --watch | Auto-restart the process when it exits |
| --max-restarts N | Limit restarts before giving up (default 10) |
| --env-file FILE | Load KEY=VALUE pairs from a file |
| --name NAME | Custom process name (default: binary basename) |
| --cwd DIR | Working directory for the child process |

## Features

### Auto-restart (--watch)

Run a process as a supervised service. If it crashes, superbg re-spawns it
with exponential backoff (1s, 2s, 4s, 8s, 16s, 30s max).

```bash
superbg run --watch --max-restarts=5 ./server
```

superbg stays alive as a monitor process. It forwards SIGTERM from
`superbg stop` to the child for graceful shutdown. State is tracked in
`~/.superbg/state.json` so restarts are visible via `list` and `status`.

### Graceful stop with timeout

```bash
superbg stop --timeout 15 1
```

Sends SIGTERM and waits up to 15 seconds. If the process is still alive,
sends SIGKILL. Default timeout is 5 seconds.

### Environment file support

```bash
superbg run --env-file .env ./app
```

Parses KEY=VALUE lines (supports # comments and optional export prefix)
and injects them into the child process environment.

### Custom name and working directory

```bash
superbg run --name my-api --cwd /opt/myapp ./bin/server
```

### Log rotation

Logs are automatically trimmed to prevent disk exhaustion. Default limits:
- 2000 lines, or
- 1MB file size

Trimming happens when viewing logs with `superbg logs` and after each
restart in --watch mode.

### Crash-loop detection

If a process crashes 3 times in a row with under 1 second of runtime each,
superbg prints a warning to stderr so you know something is broken.

### JSON output

```bash
superbg list --json
superbg status --json 1
superbg logs --json 1
```

All three commands support `--json` for machine-readable output, useful for
scripts and monitoring tools.

### Shell completions

```bash
superbg completion bash > /etc/bash_completion.d/superbg
superbg completion zsh > /usr/local/share/zsh/site-functions/_superbg
superbg completion fish > ~/.config/fish/completions/superbg.fish
```

## How it works

1. superbg run spawns the command in a new session (setsid), fully detached
   from the terminal -- no SIGHUP, no accidental Ctrl+C.
2. Stdout and stderr are captured to ~/.superbg/logs/<id>.log.
3. Process PID and metadata are saved to ~/.superbg/state.json.
4. With --watch, superbg remains as a monitor, re-spawning the child on
   exit and trimming logs.
5. All commands read from the state file, so processes survive reboots and
   terminal closures.

## Install

```
go install github.com/javimosch/superbg@latest
```

Or download a binary from the releases page at
https://github.com/javimosch/superbg/releases.

## Requirements

- Linux (uses setsid syscall)
- Go 1.22+ (to build)
