package cmd

import (
	"fmt"

	"github.com/lordofthemind/gollama/services"
	"github.com/spf13/cobra"
)

var (
	Config     string
	ConfigPath string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gollama",
	Short: "A CLI wrapper for Ollama",
	Long:  "Gollama is a CLI tool that ensures Ollama is properly installed, running, and configured.",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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
		Config, ConfigPath, err := services.LoadConfig()
		if err != nil {
			fmt.Println("Error loading configuration:", err)
			return
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gollama.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
