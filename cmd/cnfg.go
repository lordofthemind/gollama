package cmd

import (
	"fmt"

	"github.com/lordofthemind/gollama/helpers"
	"github.com/lordofthemind/gollama/services"
	"github.com/spf13/cobra"
)

var (
	tempFlag   float64
	pModelFlag string
	sModelFlag string
	tModelFlag string
)

// cnfgCmd represents the configuration management command
var cnfgCmd = &cobra.Command{
	Use:   "cnfg",
	Short: "Manage the application configuration",
	Long:  "The cnfg command allows you to manage Gollama's configuration, including models and temperature settings.",
	Run: func(cmd *cobra.Command, args []string) {
		// Load config and path
		config, configPath, err := services.LoadConfig()
		if err != nil {
			fmt.Println("Error loading configuration:", err)
			return
		}

		if !config.SetupCompleted {
			// Loop for configuring models and temperatures recursively
			for {
				// Prompt user to select models and temperatures
				services.SelectModelWithTemperature(&config)

				// Display the current configuration to confirm
				services.DisplayConfig(config)

				// Confirm with the user
				if helpers.ConfirmAction("Do you want to confirm this configuration? (y/n): ") {
					// Save config if confirmed
					config.SetupCompleted = true
					if err := services.SaveConfig(config, configPath); err != nil {
						fmt.Println("Error saving configuration:", err)
					} else {
						fmt.Println("Configuration updated successfully.")
					}
					break // Exit the loop if confirmed
				} else {
					fmt.Println("Configuration changes discarded. Please select again.")
				}
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(cnfgCmd)

	// Define flags for the command
	cnfgCmd.Flags().Float64VarP(&tempFlag, "temp", "t", 0.5, "Set the model temperature (0.1-1.0)")
	cnfgCmd.Flags().StringVarP(&pModelFlag, "primary", "p", "", "Set the primary model")
	cnfgCmd.Flags().StringVarP(&sModelFlag, "secondary", "s", "", "Set the secondary model")
	cnfgCmd.Flags().StringVarP(&tModelFlag, "tertiary", "e", "", "Set the tertiary model")
}
