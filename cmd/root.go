package cmd

import (
	"fmt"
	"os"

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
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is specified, default to running the chat command
		if len(args) == 0 {
			chatCmd.Run(cmd, args)
		} else {
			fmt.Println("Command not recognized. Use gollama --help for a list of commands.")
		}
	},
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
	rootCmd.AddCommand(chatCmd) // Add chat command as a subcommand
	rootCmd.AddCommand(cnfgCmd) // Add other commands here
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(updtCmd)

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		// Ensure Ollama is installed
		if !utils.CheckOllamaInstallation() {
			os.Exit(1)
		}

		// Ensure Ollama is running
		if !utils.CheckAndStartOllama() {
			os.Exit(1)
		}

		// Ensure gollama.yaml exists and load configuration
		config, configPath, err := configs.LoadGlobalConfig()
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
