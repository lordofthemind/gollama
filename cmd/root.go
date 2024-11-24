package cmd

import (
	"fmt"

	"github.com/lordofthemind/gollama/configs"
	"github.com/lordofthemind/gollama/services"
	"github.com/spf13/cobra"
)

var (
	Config     *configs.GollamaGlobalConfig // Pointer to hold the global configuration
	ConfigPath string
)

var rootCmd = &cobra.Command{
	Use:     "gollama",             // Main command name
	Aliases: []string{"gl", "glm"}, // Alias for the main command (so 'gl' can be used instead of 'gollama')
	Short:   "A CLI wrapper for Ollama",
	Long:    "Gollama is a CLI tool that ensures Ollama is properly installed, running, and configured.",
	Run: func(cmd *cobra.Command, args []string) {
		// Default action when no subcommand is passed, invoke 'chatCmd'
		// fmt.Println("Running the default action (chat)...")
		// Directly call the chatCmd's Run method
		chatCmd.Run(cmd, args)
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Pre-run checks for all commands
		if !services.IsOllamaInstalled() {
			fmt.Println("Exiting due to missing Ollama installation.")
			return
		}
		if !services.EnsureOllamaIsRunning() {
			fmt.Println("Exiting due to failure in starting Ollama.")
			return
		}
		// Load configuration and path
		var err error
		var loadedConfig configs.GollamaGlobalConfig
		loadedConfig, ConfigPath, err = services.LoadConfig()
		if err != nil {
			fmt.Println("Error loading configuration:", err)
			return
		}
		Config = &loadedConfig // Convert the struct to a pointer

		// Proceed with recursive configuration setup if SetupCompleted is false
		if !Config.SetupCompleted {
			startConfigurationSetup(Config, ConfigPath)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func init() {
	// Initialize persistent flags or other settings for rootCmd if needed
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Add other subcommands to the root command
	// Example: rootCmd.AddCommand(updtCmd)  // Add your subcommands here
}
