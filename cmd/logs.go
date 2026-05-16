package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/superbg/cli/state"
)

func Logs(idOrPID string, follow bool) error {
	s, err := state.Load()
	if err != nil {
		return err
	}

	job := s.FindByIDOrPID(idOrPID)
	if job == nil {
		return fmt.Errorf("no process found for: %s", idOrPID)
	}

	logFile, err := state.LogFile(job.ID)
	if err != nil {
		return err
	}

	f, err := os.Open(logFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No log output yet.")
			return nil
		}
		return fmt.Errorf("open log: %w", err)
	}
	defer f.Close()

	if follow {
		return followLog(f)
	}

	_, err = io.Copy(os.Stdout, f)
	return err
}

func followLog(f *os.File) error {
	_, err := io.Copy(os.Stdout, f)
	if err != nil {
		return err
	}

	info, err := f.Stat()
	if err != nil {
		return err
	}
	offset := info.Size()

	buf := make([]byte, 4096)
	for {
		n, err := f.ReadAt(buf, offset)
		if n > 0 {
			os.Stdout.Write(buf[:n])
			offset += int64(n)
		}
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			info, err := f.Stat()
			if err != nil {
				return err
			}
			if info.Size() < offset {
				offset = 0
			}
		}
	}
}
