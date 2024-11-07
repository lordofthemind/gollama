package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"

	"github.com/lordofthemind/gollama/configs"
	"github.com/lordofthemind/gollama/utils"
)

var rootCmd = &cobra.Command{
	Use:   "gollama",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		// Ensure Ollama is installed
		if !utils.CheckOllamaInstallation() {
			os.Exit(1)
		}

		// Ensure Ollama is running
		if !CheckAndStartOllama() {
			os.Exit(1)
		}

		// Ensure gollama.yaml exists and load configuration
		_, configPath, err := configs.LoadGlobalConfig()
		if err != nil {
			fmt.Println("Error loading configuration:", err)
			os.Exit(1)
		}

		// If configuration setup is incomplete, initiate setup
		config, _, err := configs.LoadGlobalConfig()
		if err != nil {
			fmt.Println("Error loading configuration:", err)
			os.Exit(1)
		}

		if !config.SetupCompleted {
			fmt.Println("Configuration setup required. Starting setup...")
			initiateGlobalConfigurationSetup(&config, configPath)
		}
	}
}

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
