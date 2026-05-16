package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type JobStatus string

const (
	StatusRunning JobStatus = "running"
	StatusStopped JobStatus = "stopped"
	StatusKilled  JobStatus = "killed"
	StatusExited  JobStatus = "exited"
)

type Job struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Command   []string  `json:"command"`
	PID       int       `json:"pid"`
	Status    JobStatus `json:"status"`
	StartedAt string    `json:"started_at"`
	StoppedAt string    `json:"stopped_at,omitempty"`
	ExitCode  int       `json:"exit_code,omitempty"`
}

type State struct {
	Jobs   []Job `json:"jobs"`
	NextID int   `json:"next_id"`
}

func homedir() (string, error) {
	h := os.Getenv("HOME")
	if h == "" {
		return "", fmt.Errorf("HOME not set")
	}
	return h, nil
}

func baseDir() (string, error) {
	h, err := homedir()
	if err != nil {
		return "", err
	}
	return filepath.Join(h, ".superbg"), nil
}

func StateFile() (string, error) {
	d, err := baseDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, "state.json"), nil
}

func LogDir() (string, error) {
	d, err := baseDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(d, "logs"), nil
}

func LogFile(id int) (string, error) {
	dir, err := LogDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, fmt.Sprintf("%d.log", id)), nil
}

func InitDirs() error {
	d, err := baseDir()
	if err != nil {
		return err
	}
	logdir, err := LogDir()
	if err != nil {
		return err
	}
	for _, p := range []string{d, logdir} {
		if err := os.MkdirAll(p, 0755); err != nil {
			return fmt.Errorf("create dir %s: %w", p, err)
		}
	}
	return nil
}

func Load() (*State, error) {
	s := &State{NextID: 1}
	p, err := StateFile()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return s, nil
		}
		return nil, fmt.Errorf("read state: %w", err)
	}
	if err := json.Unmarshal(data, s); err != nil {
		return nil, fmt.Errorf("parse state: %w", err)
	}
	return s, nil
}

func (s *State) Save() error {
	p, err := StateFile()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("encode state: %w", err)
	}
	if err := os.WriteFile(p, data, 0644); err != nil {
		return fmt.Errorf("write state: %w", err)
	}
	return nil
}

func (s *State) AddJob(name string, command []string, pid int) Job {
	job := Job{
		ID:        s.NextID,
		Name:      name,
		Command:   command,
		PID:       pid,
		Status:    StatusRunning,
		StartedAt: time.Now().UTC().Format(time.RFC3339),
	}
	s.NextID++
	s.Jobs = append(s.Jobs, job)
	return job
}

func (s *State) RemoveJob(id int) {
	idx := -1
	for i, j := range s.Jobs {
		if j.ID == id {
			idx = i
			break
		}
	}
	if idx >= 0 {
		s.Jobs = append(s.Jobs[:idx], s.Jobs[idx+1:]...)
	}
}

func (s *State) FindByID(idStr string) *Job {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil
	}
	for i := range s.Jobs {
		if s.Jobs[i].ID == id {
			return &s.Jobs[i]
		}
	}
	return nil
}

func (s *State) FindByPID(pid int) *Job {
	for i := range s.Jobs {
		if s.Jobs[i].PID == pid {
			return &s.Jobs[i]
		}
	}
	return nil
}

func (s *State) FindByIDOrPID(idOrPID string) *Job {
	if job := s.FindByID(idOrPID); job != nil {
		return job
	}
	pid, err := strconv.Atoi(idOrPID)
	if err != nil {
		return nil
	}
	return s.FindByPID(pid)
}
