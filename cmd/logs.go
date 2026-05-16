package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/superbg/cli/state"
)

const (
	maxLogBytes = 1 * 1024 * 1024
	maxLogLines = 2000
	trimLines   = 1000
)

func trimLog(path string) {
	info, err := os.Stat(path)
	if err != nil || info.Size() == 0 {
		return
	}

	if info.Size() > maxLogBytes {
		sz := info.Size()
		data, _ := os.ReadFile(path)
		if data == nil {
			return
		}
		keep := int(sz / 2)
		start := len(data) - keep
		for start < len(data) && data[start] != '\n' {
			start++
		}
		if start < len(data) {
			start++
		}
		if start < len(data) {
			os.WriteFile(path, data[start:], 0644)
		}
		return
	}

	data, _ := os.ReadFile(path)
	if data == nil {
		return
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) <= maxLogLines {
		return
	}
	start := len(lines) - trimLines
	if start < 0 {
		start = 0
	}
	os.WriteFile(path, []byte(strings.Join(lines[start:], "\n")), 0644)
}

func Logs(idOrPID string, follow bool, asJSON bool) error {
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

	trimLog(logFile)

	f, err := os.Open(logFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No log output yet.")
			return nil
		}
		return fmt.Errorf("open log: %w", err)
	}
	defer f.Close()

	if asJSON {
		data, _ := os.ReadFile(logFile)
		out := map[string]interface{}{
			"id":    job.ID,
			"name":  job.Name,
			"pid":   job.PID,
			"logs":  string(data),
			"lines": len(strings.Split(string(data), "\n")) - 1,
		}
		return json.NewEncoder(os.Stdout).Encode(out)
	}

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
