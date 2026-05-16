package cmd

import (
	"fmt"
	"time"

	"github.com/superbg/cli/process"
	"github.com/superbg/cli/state"
)

func Kill(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: superbg kill <id|pid>")
	}

	target := args[len(args)-1]

	s, err := state.Load()
	if err != nil {
		return err
	}

	job := s.FindByIDOrPID(target)
	if job == nil {
		return fmt.Errorf("no process found for: %s", target)
	}

	killPID := job.PID
	alsoKill := []int{}

	if job.MonitorPID > 0 {
		alsoKill = append(alsoKill, job.MonitorPID)
	}

	if !process.IsRunning(killPID) {
		job.Status = state.StatusExited
		s.Save()
		return fmt.Errorf("process %d is not running", killPID)
	}

	fmt.Printf("Killing process %d (%s)...\n", killPID, job.Name)

	if err := process.Kill(killPID); err != nil {
		return fmt.Errorf("kill %d: %w", killPID, err)
	}

	for _, pid := range alsoKill {
		if process.IsRunning(pid) {
			process.Kill(pid)
		}
	}

	time.Sleep(100 * time.Millisecond)
	for _, pid := range alsoKill {
		if process.IsRunning(pid) {
			process.Kill(pid)
		}
	}

	job.Status = state.StatusKilled
	job.StoppedAt = currentTimestamp()
	job.PID = 0

	if err := s.Save(); err != nil {
		return err
	}

	fmt.Printf("Process %d killed.\n", killPID)
	return nil
}
