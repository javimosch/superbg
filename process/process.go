package process

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func Run(command []string, logFile string, extraEnv []string) (*exec.Cmd, error) {
	if len(command) == 0 {
		return nil, fmt.Errorf("no command specified")
	}

	cmd := exec.Command(command[0], command[1:]...)

	f, err := os.Create(logFile)
	if err != nil {
		return nil, fmt.Errorf("create log file: %w", err)
	}

	cmd.Stdout = f
	cmd.Stderr = f
	cmd.Stdin = nil

	if len(extraEnv) > 0 {
		cmd.Env = append(os.Environ(), extraEnv...)
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	if err := cmd.Start(); err != nil {
		f.Close()
		return nil, fmt.Errorf("start process: %w", err)
	}

	return cmd, nil
}

func Signal(pid int, sig os.Signal) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("find process %d: %w", pid, err)
	}
	return p.Signal(sig)
}

func Stop(pid int) error {
	return Signal(pid, syscall.SIGTERM)
}

func Kill(pid int) error {
	return Signal(pid, syscall.SIGKILL)
}

func IsRunning(pid int) bool {
	err := Signal(pid, syscall.Signal(0))
	return err == nil
}
