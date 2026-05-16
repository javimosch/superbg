package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/superbg/cli/process"
	"github.com/superbg/cli/state"
)

func Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: superbg run <command> [args...]")
	}

	if err := state.InitDirs(); err != nil {
		return fmt.Errorf("init dirs: %w", err)
	}

	s, err := state.Load()
	if err != nil {
		return err
	}

	name := filepath.Base(args[0])

	logFile, err := state.LogFile(s.NextID)
	if err != nil {
		return err
	}

	cmd, err := process.Run(args, logFile)
	if err != nil {
		return err
	}

	pid := cmd.Process.Pid
	job := s.AddJob(name, args, pid)

	if err := s.Save(); err != nil {
		return err
	}

	fmt.Printf("[%d] %d\n", job.ID, pid)
	fmt.Printf("Logs: %s\n", logFile)

	go func() {
		cmd.Wait()
	}()

	return nil
}
