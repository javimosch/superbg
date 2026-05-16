package cmd

import "time"

func currentTimestamp() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05 UTC")
}

func formatCommand(cmd []string) string {
	if len(cmd) == 0 {
		return ""
	}
	result := cmd[0]
	for _, a := range cmd[1:] {
		result += " " + a
	}
	return result
}
