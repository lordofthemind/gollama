package helpers

import (
	"os/exec"
)

// Check if a command exists in the system PATH
func IsCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// Check if a command can execute successfully
func IsCommandRunning(command string, args ...string) bool {
	cmd := exec.Command(command, args...)
	err := cmd.Run()
	return err == nil
}
