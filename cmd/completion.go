package cmd

import (
	"fmt"
	"os"
)

func Completion(shell string) error {
	switch shell {
	case "bash":
		return printBash()
	case "zsh":
		return printZsh()
	case "fish":
		return printFish()
	default:
		return fmt.Errorf("unsupported shell: %s (use: bash, zsh, fish)", shell)
	}
}

func printBash() error {
	_, err := os.Stdout.WriteString(`_superbg() {
    local cur prev opts cmd
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    opts="run list stop kill logs status attach help clean rm completion"

    if [[ ${COMP_CWORD} -eq 1 ]]; then
        COMPREPLY=($(compgen -W "${opts}" -- ${cur}))
        return 0
    fi

    cmd="${COMP_WORDS[1]}"
    case "${cmd}" in
        run)
            COMPREPLY=($(compgen -W "--watch --max-restarts --env-file --name --cwd" -- ${cur}))
            ;;
        stop)
            COMPREPLY=($(compgen -W "--timeout" -- ${cur}))
            ;;
        logs)
            COMPREPLY=($(compgen -W "--follow -f" -- ${cur}))
            ;;
        completion)
            COMPREPLY=($(compgen -W "bash zsh fish" -- ${cur}))
            ;;
        status|kill|rm|attach)
            COMPREPLY=()
            ;;
    esac
    return 0
}
complete -F _superbg superbg
`)
	return err
}

func printZsh() error {
	_, err := os.Stdout.WriteString(`#compdef superbg
_superbg() {
    local -a opts
    opts=(
        'run:Run a command in the background'
        'list:List background processes'
        'stop:Stop a process'
        'kill:Kill a process'
        'logs:Show process logs'
        'status:Show process status'
        'attach:Follow process logs in real-time'
        'help:Show help'
        'clean:Remove completed processes'
        'rm:Remove a process from tracking'
        'completion:Generate shell completion'
    )
    _describe 'superbg' opts
}
_superbg "$@"
`)
	return err
}

func printFish() error {
	_, err := os.Stdout.WriteString(`complete -c superbg -f
complete -c superbg -n '__fish_use_subcommand' -a run -d 'Run a command in the background'
complete -c superbg -n '__fish_use_subcommand' -a list -d 'List background processes'
complete -c superbg -n '__fish_use_subcommand' -a stop -d 'Stop a process'
complete -c superbg -n '__fish_use_subcommand' -a kill -d 'Kill a process'
complete -c superbg -n '__fish_use_subcommand' -a logs -d 'Show process logs'
complete -c superbg -n '__fish_use_subcommand' -a status -d 'Show process status'
complete -c superbg -n '__fish_use_subcommand' -a attach -d 'Follow process logs'
complete -c superbg -n '__fish_use_subcommand' -a clean -d 'Remove completed processes'
complete -c superbg -n '__fish_use_subcommand' -a rm -d 'Remove a process from tracking'
complete -c superbg -n '__fish_use_subcommand' -a completion -d 'Generate shell completion'

complete -c superbg -n '__fish_seen_subcommand_from run' -l watch -d 'Auto-restart on crash'
complete -c superbg -n '__fish_seen_subcommand_from run' -l max-restarts -d 'Max restarts'
complete -c superbg -n '__fish_seen_subcommand_from run' -l env-file -d 'Load env file'
complete -c superbg -n '__fish_seen_subcommand_from run' -l name -d 'Custom process name'
complete -c superbg -n '__fish_seen_subcommand_from run' -l cwd -d 'Working directory'
complete -c superbg -n '__fish_seen_subcommand_from stop' -l timeout -d 'Graceful stop timeout'
complete -c superbg -n '__fish_seen_subcommand_from logs' -s f -l follow -d 'Follow log output'
`)
	return err
}
