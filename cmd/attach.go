package cmd

import (
	"fmt"
	"os"

	"github.com/superbg/cli/process"
	"github.com/superbg/cli/state"
)

func Attach(idOrPID string) error {
	s, err := state.Load()
	if err != nil {
		return err
	}

	job := s.FindByIDOrPID(idOrPID)
	if job == nil {
		return fmt.Errorf("no process found for: %s", idOrPID)
	}

	if !process.IsRunning(job.PID) {
		return fmt.Errorf("process %d (%s) is not running", job.PID, job.Name)
	}

	logFile, err := state.LogFile(job.ID)
	if err != nil {
		return err
	}

	f, err := os.Open(logFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No log output yet. Waiting for output...")
			return nil
		}
		return fmt.Errorf("open log: %w", err)
	}
	defer f.Close()

	return followLog(f)
}
