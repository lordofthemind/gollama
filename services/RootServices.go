package services

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/lordofthemind/gollama/helpers"
)

// Check if Ollama is installed
func IsOllamaInstalled() bool {
	if !helpers.IsCommandAvailable("ollama") {
		fmt.Println("Ollama is not installed. Please install Ollama to proceed with Gollama setup.")
		return false
	}
	return true
}

// Check if Ollama is running, and start it if not
func EnsureOllamaIsRunning() bool {
	if !helpers.IsCommandRunning("ollama", "ps") {
		fmt.Println("Ollama is not running. Attempting to start Ollama...")

		// Start Ollama (replace 'ollama serve' with the appropriate command if different on Windows)
		cmd := exec.Command("ollama", "serve")
		err := cmd.Start()
		if err != nil {
			fmt.Println("Failed to start Ollama:", err)
			return false
		}

		// Wait for Ollama to initialize
		time.Sleep(2 * time.Second)

		// Verify if Ollama started successfully
		if !helpers.IsCommandRunning("ollama", "ps") {
			fmt.Println("Ollama failed to start. Please start it manually and try again.")
			return false
		}
		fmt.Println("Ollama started successfully.")
	}
	return true
}
