package cmd

import (
	"fmt"

	"github.com/superbg/cli/process"
	"github.com/superbg/cli/state"
)

func Status(idOrPID string) error {
	s, err := state.Load()
	if err != nil {
		return err
	}

	job := s.FindByIDOrPID(idOrPID)
	if job == nil {
		return fmt.Errorf("no process found for: %s", idOrPID)
	}

	alive := process.IsRunning(job.PID)
	status := string(job.Status)

	if job.Status == state.StatusRunning && !alive {
		status = "exited (unexpected)"
	}

	fmt.Printf("ID:        %d\n", job.ID)
	fmt.Printf("Name:      %s\n", job.Name)
	fmt.Printf("Command:   %s\n", formatCommand(job.Command))
	fmt.Printf("PID:       %d\n", job.PID)
	fmt.Printf("Status:    %s\n", status)
	fmt.Printf("Started:   %s\n", job.StartedAt)
	if job.StoppedAt != "" {
		fmt.Printf("Stopped:   %s\n", job.StoppedAt)
	}
	if job.ExitCode != 0 {
		fmt.Printf("Exit Code: %d\n", job.ExitCode)
	}

	logFile, err := state.LogFile(job.ID)
	if err == nil {
		fmt.Printf("Logs:      %s\n", logFile)
	}

	return nil
}
