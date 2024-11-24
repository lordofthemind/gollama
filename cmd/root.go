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

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gollama",
	Short: "A CLI wrapper for Ollama",
	Long:  "Gollama is a CLI tool that ensures Ollama is properly installed, running, and configured.",
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
}
