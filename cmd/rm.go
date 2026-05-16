package cmd

import (
	"fmt"

	"github.com/superbg/cli/state"
)

func Rm(idOrPID string) error {
	s, err := state.Load()
	if err != nil {
		return err
	}

	job := s.FindByIDOrPID(idOrPID)
	if job == nil {
		return fmt.Errorf("no process found for: %s", idOrPID)
	}

	id := job.ID
	s.RemoveJob(id)

	if err := s.Save(); err != nil {
		return err
	}

	fmt.Printf("Removed process %d from tracking.\n", id)
	return nil
}
