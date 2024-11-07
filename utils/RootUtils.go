package utils

import (
	"fmt"
	"os/exec"
	"time"
)

// Function to check if Ollama is running, and start it if not
func CheckAndStartOllama() bool {
	// Check if Ollama is running by executing `ollama list`
	if !isOllamaRunning() {
		fmt.Println("Ollama is not running. Attempting to start Ollama...")

		cmd := exec.Command("ollama", "list") // Adjust the command to start Ollama if necessary
		err := cmd.Start()
		if err != nil {
			fmt.Println("Failed to start Ollama:", err)
			return false
		}

		// Wait for Ollama to start
		time.Sleep(2 * time.Second)

		// Verify Ollama started successfully
		if !isOllamaRunning() {
			fmt.Println("Ollama failed to start. Please start Ollama manually and try again.")
			return false
		}
		fmt.Println("Ollama started successfully.")
	}
	return true
}

// Function to check if Ollama is running by attempting to execute `ollama list`
func isOllamaRunning() bool {
	cmd := exec.Command("ollama", "list")
	if err := cmd.Run(); err != nil {
		fmt.Println("Ollama is not currently running.")
		return false
	}
	return true
}
