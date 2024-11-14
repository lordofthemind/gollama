package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/lordofthemind/gollama/utils"

)

var logger utils.Logger

var rootCmd = &cobra.Command{
	Use:   "gollama",
	Short: "A brief description of your application",
	Long:  `A longer description for the gollama CLI tool with example usage.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Run default command if none specified
		if len(args) == 0 {
			// chatCmd.Run(cmd, args)
		} else {
			fmt.Println("Command not recognized. Use gollama --help for a list of commands.")
		}
	},
}

func Execute() {
	logFilePath := filepath.Join(os.Getenv("HOME"), ".config", "Gollama", "gollama.log")
	var err error
	logger, err = utils.NewLogger(logFilePath)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	if err := rootCmd.Execute(); err != nil {
		logger.LogFatal("Failed to execute root command:", err)
		os.Exit(1)
	}
}

func init() {
	// Add subcommands to rootCmd here
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		utils.SetupEnvironment(logger)
	}
}
