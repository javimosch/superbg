package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/superbg/cli/process"
	"github.com/superbg/cli/state"
)

func Stop(args []string) error {
	timeout := 5
	target := ""

	for _, a := range args {
		if a == "--timeout" || a == "-t" {
			continue
		}
		if strings.HasPrefix(a, "--timeout=") {
			if n, err := strconv.Atoi(strings.TrimPrefix(a, "--timeout=")); err == nil {
				timeout = n
			}
			continue
		}
		// Check if previous arg was --timeout
		target = a
	}

	// Re-scan for --timeout N pattern
	for i := 0; i < len(args); i++ {
		if args[i] == "--timeout" || args[i] == "-t" {
			if i+1 < len(args) {
				if n, err := strconv.Atoi(args[i+1]); err == nil {
					timeout = n
				}
			}
		}
	}

	if target == "" {
		return fmt.Errorf("usage: superbg stop [--timeout N] <id|pid>")
	}

	s, err := state.Load()
	if err != nil {
		return err
	}

	job := s.FindByIDOrPID(target)
	if job == nil {
		return fmt.Errorf("no process found for: %s", target)
	}

	targetPID := job.PID
	if job.MonitorPID > 0 {
		targetPID = job.MonitorPID
	}

	if !process.IsRunning(targetPID) {
		if job.MonitorPID > 0 && job.PID > 0 && !process.IsRunning(job.PID) {
			job.Status = state.StatusExited
			s.Save()
			return fmt.Errorf("process is not running")
		}
		job.Status = state.StatusExited
		s.Save()
		return fmt.Errorf("process %d is not running", targetPID)
	}

	fmt.Printf("Stopping process %d (%s)...\n", targetPID, job.Name)

	if err := process.Stop(targetPID); err != nil {
		return fmt.Errorf("stop %d: %w", targetPID, err)
	}

	waited := 0
	for waited < timeout {
		if !process.IsRunning(targetPID) {
			break
		}
		time.Sleep(time.Second)
		waited++
	}

	if process.IsRunning(targetPID) {
		fmt.Printf("Process %d did not stop after %ds, sending SIGKILL...\n", targetPID, timeout)
		if err := process.Kill(targetPID); err != nil {
			return fmt.Errorf("kill %d: %w", targetPID, err)
		}
	}

	job.Status = state.StatusStopped
	job.StoppedAt = currentTimestamp()
	if job.MonitorPID > 0 {
		job.PID = 0
	}

	if err := s.Save(); err != nil {
		return err
	}

	fmt.Printf("Process %d stopped.\n", targetPID)
	return nil
}
