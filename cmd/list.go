package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/superbg/cli/process"
	"github.com/superbg/cli/state"
)

func List() error {
	s, err := state.Load()
	if err != nil {
		return err
	}

	if len(s.Jobs) == 0 {
		fmt.Println("No background processes.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tSTATUS\tPID\tNAME\tCOMMAND")

	for _, job := range s.Jobs {
		status := string(job.Status)
		pid := job.PID

		if job.Status == state.StatusRunning {
			if process.IsRunning(job.PID) {
				status = "running"
			} else {
				status = "exited"
				pid = 0
			}
		}

		fmt.Fprintf(w, "%d\t%s\t%d\t%s\t%s\n",
			job.ID, status, pid, job.Name, formatCommand(job.Command))
	}

	w.Flush()
	return nil
}
