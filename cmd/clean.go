package cmd

import (
	"fmt"

	"github.com/superbg/cli/state"
)

func Clean() error {
	s, err := state.Load()
	if err != nil {
		return err
	}

	before := len(s.Jobs)
	var kept []state.Job
	for _, j := range s.Jobs {
		if j.Status == state.StatusRunning {
			kept = append(kept, j)
		}
	}
	s.Jobs = kept

	if err := s.Save(); err != nil {
		return err
	}

	removed := before - len(kept)
	if removed == 0 {
		fmt.Println("No completed processes to clean.")
	} else {
		fmt.Printf("Removed %d completed process(es). %d remaining.\n", removed, len(kept))
	}
	return nil
}
