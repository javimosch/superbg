package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/superbg/cli/process"
	"github.com/superbg/cli/state"
)

func Status(idOrPID string, asJSON bool) error {
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

	if asJSON {
		info := map[string]interface{}{
			"id":            job.ID,
			"name":          job.Name,
			"command":       job.Command,
			"pid":           job.PID,
			"status":        status,
			"started_at":    job.StartedAt,
			"stopped_at":    job.StoppedAt,
			"exit_code":     job.ExitCode,
			"auto_restart":  job.AutoRestart,
			"restart_count": job.RestartCount,
			"max_restarts":  job.MaxRestarts,
			"monitor_pid":   job.MonitorPID,
		}
		if job.AutoRestart && job.MonitorPID > 0 {
			info["monitor_alive"] = process.IsRunning(job.MonitorPID)
		}
		logFile, err := state.LogFile(job.ID)
		if err == nil {
			info["log_file"] = logFile
		}
		return json.NewEncoder(os.Stdout).Encode(info)
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
	if job.AutoRestart {
		fmt.Printf("Watch:     yes (restarts %d/%d)\n", job.RestartCount, job.MaxRestarts)
		if job.MonitorPID > 0 {
			monitorAlive := process.IsRunning(job.MonitorPID)
			aliveStr := "alive"
			if !monitorAlive {
				aliveStr = "dead"
			}
			fmt.Printf("Monitor:   %d (%s)\n", job.MonitorPID, aliveStr)
		}
	}

	logFile, err := state.LogFile(job.ID)
	if err == nil {
		fmt.Printf("Logs:      %s\n", logFile)
	}

	return nil
}
