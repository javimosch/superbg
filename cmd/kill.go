package cmd

import (
	"fmt"

	"github.com/superbg/cli/process"
	"github.com/superbg/cli/state"
)

func Kill(idOrPID string) error {
	s, err := state.Load()
	if err != nil {
		return err
	}

	job := s.FindByIDOrPID(idOrPID)
	if job == nil {
		return fmt.Errorf("no process found for: %s", idOrPID)
	}

	if !process.IsRunning(job.PID) {
		job.Status = state.StatusExited
		s.Save()
		return fmt.Errorf("process %d is not running", job.PID)
	}

	fmt.Printf("Killing process %d (%s)...\n", job.PID, job.Name)

	if err := process.Kill(job.PID); err != nil {
		return fmt.Errorf("kill %d: %w", job.PID, err)
	}

	job.Status = state.StatusKilled
	job.StoppedAt = currentTimestamp()

	if err := s.Save(); err != nil {
		return err
	}

	fmt.Printf("Process %d killed.\n", job.PID)
	return nil
}
