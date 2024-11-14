package utils

import (
	"os"
	"os/exec"

	"github.com/lordofthemind/gollama/configs"
)

// SetupEnvironment ensures all necessary setup steps are complete
func SetupEnvironment(logger Logger) {
	// Step 1: Ensure Ollama is installed
	if !CheckOllamaInstallation(logger) {
		logger.LogFatal("Ollama is not installed. Please install it to proceed.")
		os.Exit(1)
	}

	// Step 2: Ensure Ollama is running
	if !CheckAndStartOllama(logger) {
		logger.LogFatal("Failed to start Ollama. Please start it manually and try again.")
		os.Exit(1)
	}

	// Step 3: Ensure configuration is loaded
	config, configPath, err := configs.LoadGlobalConfig()
	if err != nil {
		logger.LogFatal("Error loading configuration:", err)
		os.Exit(1)
	}

	// Step 4: Initiate setup if configuration is incomplete
	if !config.SetupCompleted {
		logger.LogInfo("Configuration setup required. Initiating setup...")
		initiateGlobalConfigurationSetup(logger, &config, configPath)
	}
}

// CheckOllamaInstallation verifies that Ollama is installed by checking its presence in PATH.
func CheckOllamaInstallation(logger Logger) bool {
	_, err := exec.LookPath("ollama")
	if err != nil {
		logger.LogError("Ollama is not installed or not in PATH.")
		return false
	}
	logger.LogInfo("Ollama installation verified.")
	return true
}

// CheckAndStartOllama ensures Ollama is running or attempts to start it.
func CheckAndStartOllama(logger Logger) bool {
	// Check if Ollama is running by attempting a basic command

	err := exec.Command("ollama").Run()
	if err == nil {
		logger.LogInfo("Ollama is running.")
		return true
	}

	// If not running, attempt to start Ollama
	logger.LogInfo("Ollama is not running. Attempting to start...")
	startCmd := exec.Command("ollama", "list") // Adjust the command as per Ollama's actual start command
	if err := startCmd.Start(); err != nil {
		logger.LogError("Failed to start Ollama:", err)
		return false
	}

	// Wait for Ollama to start
	err = startCmd.Wait()
	if err != nil {
		logger.LogError("Ollama failed to start:", err)
		return false
	}
	logger.LogInfo("Ollama started successfully.")
	return true
}

// initiateGlobalConfigurationSetup starts the initial configuration setup
func initiateGlobalConfigurationSetup(logger Logger, config *configs.GollamaGlobalConfig, configPath string) {
	// Custom configuration setup logic here
	// After setup, mark as completed and save the configuration
	config.SetupCompleted = true
	if err := configs.SaveGlobalConfig(*config, configPath); err != nil {
		logger.LogFatal("Failed to save configuration during setup:", err)
		os.Exit(1)
	}
	logger.LogInfo("Configuration setup completed successfully.")
}
