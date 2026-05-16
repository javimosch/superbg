package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/superbg/cli/process"
	"github.com/superbg/cli/state"
)

type runFlags struct {
	watch      bool
	maxRestarts int
	envFile    string
	cmdArgs    []string
}

func parseRunFlags(args []string) (*runFlags, error) {
	f := &runFlags{maxRestarts: 10}
	i := 0
	for i < len(args) {
		a := args[i]
		if a == "--watch" {
			f.watch = true
			i++
		} else if a == "--max-restarts" && i+1 < len(args) {
			n, err := strconv.Atoi(args[i+1])
			if err != nil {
				return nil, fmt.Errorf("invalid --max-restarts value: %s", args[i+1])
			}
			f.maxRestarts = n
			i += 2
		} else if strings.HasPrefix(a, "--max-restarts=") {
			n, err := strconv.Atoi(strings.TrimPrefix(a, "--max-restarts="))
			if err != nil {
				return nil, fmt.Errorf("invalid --max-restarts value: %s", a)
			}
			f.maxRestarts = n
			i++
		} else if a == "--env-file" && i+1 < len(args) {
			f.envFile = args[i+1]
			i += 2
		} else if strings.HasPrefix(a, "--env-file=") {
			f.envFile = strings.TrimPrefix(a, "--env-file=")
			i++
		} else if strings.HasPrefix(a, "--") {
			return nil, fmt.Errorf("unknown flag: %s", a)
		} else {
			f.cmdArgs = args[i:]
			break
		}
	}
	if len(f.cmdArgs) == 0 {
		return nil, fmt.Errorf("usage: superbg run [--watch] [--max-restarts=N] [--env-file FILE] <command> [args...]")
	}
	return f, nil
}

func parseEnvFile(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read env file: %w", err)
	}
	var env []string
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		if strings.Contains(line, "=") {
			env = append(env, line)
		}
	}
	return env, nil
}

func Run(args []string) error {
	flags, err := parseRunFlags(args)
	if err != nil {
		return err
	}

	if err := state.InitDirs(); err != nil {
		return fmt.Errorf("init dirs: %w", err)
	}

	var extraEnv []string
	if flags.envFile != "" {
		extraEnv, err = parseEnvFile(flags.envFile)
		if err != nil {
			return err
		}
	}

	s, err := state.Load()
	if err != nil {
		return err
	}

	name := filepath.Base(flags.cmdArgs[0])

	logFile, err := state.LogFile(s.NextID)
	if err != nil {
		return err
	}

	cmd, err := process.Run(flags.cmdArgs, logFile, extraEnv)
	if err != nil {
		return err
	}

	pid := cmd.Process.Pid
	job := s.AddJob(name, flags.cmdArgs, pid)

	if flags.watch {
		monitorPID := os.Getpid()
		job.MonitorPID = monitorPID
		job.AutoRestart = true
		job.MaxRestarts = flags.maxRestarts
	}

	if err := s.Save(); err != nil {
		return err
	}

	fmt.Printf("[%d] %d\n", job.ID, pid)
	fmt.Printf("Logs: %s\n", logFile)

	if flags.watch {
		return runWatch(cmd, job.ID, flags.maxRestarts, logFile, extraEnv)
	}

	go func() {
		cmd.Wait()
	}()
	return nil
}

func backoff(attempt int) time.Duration {
	d := time.Duration(1<<min(attempt, 5)) * time.Second
	if d > 30*time.Second {
		d = 30 * time.Second
	}
	return d
}

func runWatch(initialCmd *exec.Cmd, jobID, maxRestarts int, logFile string, extraEnv []string) error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(sigCh)

	var (
		childPID   atomic.Int32
		shutdown   atomic.Bool
		done       = make(chan struct{})
	)

	childPID.Store(int32(initialCmd.Process.Pid))

	go func() {
		select {
		case sig := <-sigCh:
			shutdown.Store(true)
			pid := childPID.Load()
			if pid > 0 {
				syscall.Kill(int(pid), sig.(syscall.Signal))
			}
		case <-done:
		}
	}()
	defer close(done)

	markState := func(status state.JobStatus) {
		s, _ := state.Load()
		if job := s.FindByID(strconv.Itoa(jobID)); job != nil {
			job.Status = status
			job.StoppedAt = currentTimestamp()
			job.PID = 0
			s.Save()
		}
	}

	cmd := initialCmd
	for restarts := 0; ; {
		cmd.Wait()

		if shutdown.Load() {
			markState(state.StatusStopped)
			fmt.Println("Monitor stopped.")
			return nil
		}

		if restarts >= maxRestarts {
			markState(state.StatusExited)
			fmt.Printf("Max restarts (%d) reached. Monitor exiting.\n", maxRestarts)
			return nil
		}

		d := backoff(restarts)
		fmt.Printf("Process exited. Restarting in %v (attempt %d/%d)...\n", d, restarts+1, maxRestarts)

		select {
		case <-time.After(d):
		case <-done:
			shutdown.Store(true)
			return nil
		}

		newCmd, err := process.Run(
			initialCmd.Args,
			logFile,
			extraEnv,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Restart failed: %s\n", err)
			return err
		}

		childPID.Store(int32(newCmd.Process.Pid))
		cmd = newCmd
		restarts++

		s, _ := state.Load()
		if job := s.FindByID(strconv.Itoa(jobID)); job != nil {
			job.PID = newCmd.Process.Pid
			job.RestartCount = restarts
			job.Status = state.StatusRunning
			s.Save()
		}
		fmt.Printf("[%d] Restarted as PID %d\n", jobID, newCmd.Process.Pid)
	}
}
