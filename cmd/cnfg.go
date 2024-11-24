package cmd

import (
	"fmt"

	"github.com/lordofthemind/gollama/configs"
	"github.com/lordofthemind/gollama/helpers"
	"github.com/lordofthemind/gollama/services"
	"github.com/spf13/cobra"
)

var (
	primaryTempFlag    float64
	secondaryTempFlag  float64
	tertiaryTempFlag   float64
	primaryModelFlag   string
	secondaryModelFlag string
	tertiaryModelFlag  string
)

var cnfgCmd = &cobra.Command{
	Use:     "cnfg",
	Aliases: []string{"cfg", "config", "configuration"},
	Short:   "Manage Gollama's configuration",
	Long: `The cnfg command allows you to view and manage Gollama's configuration, including models and temperature settings.

Key Features:
- View the current configuration if no flags are provided.
- Update the primary, secondary, and tertiary models and their respective temperature settings via flags.
- Interactive configuration setup to guide users through selecting models and temperatures.
- Automatically saves changes if confirmed by the user.

Usage:
- Run 'cnfg' without flags to display the current configuration and optionally edit it.
- Use flags such as '--primary', '--primary-temp', etc., to update specific values in the configuration.`,

	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration and path
		Config, ConfigPath, err := services.LoadConfig()
		if err != nil {
			fmt.Println("Error loading configuration:", err)
			return
		}

		// If no flags are provided, display the current configuration and prompt to edit
		if !cmd.Flags().Changed("primary-temp") && !cmd.Flags().Changed("secondary-temp") &&
			!cmd.Flags().Changed("tertiary-temp") && !cmd.Flags().Changed("primary") &&
			!cmd.Flags().Changed("secondary") && !cmd.Flags().Changed("tertiary") {
			services.DisplayConfig(Config)
			if helpers.ConfirmAction("Do you want to edit the current configuration? (y/n): ") {
				startConfigurationSetup(&Config, ConfigPath)
			}
			return
		}

		// If flags are provided, update the configuration
		updated := services.UpdateConfigFromFlags(&Config, primaryTempFlag, secondaryTempFlag, tertiaryTempFlag, primaryModelFlag, secondaryModelFlag, tertiaryModelFlag)
		if updated {
			if err := services.SaveConfig(Config, ConfigPath); err != nil {
				fmt.Println("Error saving updated configuration:", err)
			} else {
				fmt.Println("Configuration updated successfully.")
			}
		} else {
			fmt.Println("No changes were made to the configuration.")
		}
	},
}

func startConfigurationSetup(config *configs.GollamaGlobalConfig, configPath string) {
	// Proceed with recursive configuration setup if SetupCompleted is false
	for {
		// Prompt user to select models and temperatures
		services.SelectModelWithTemperature(config)

		// Display the current configuration to confirm
		services.DisplayConfig(*config)

		// Confirm with the user
		if helpers.ConfirmAction("Do you want to confirm this configuration? (y/n): ") {
			// Save config if confirmed
			config.SetupCompleted = true
			if err := services.SaveConfig(*config, configPath); err != nil {
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

func init() {
	rootCmd.AddCommand(cnfgCmd)

	// Define flags for the command
	cnfgCmd.Flags().StringVarP(&primaryModelFlag, "primary", "p", "", "Set the primary model")
	cnfgCmd.Flags().StringVarP(&secondaryModelFlag, "secondary", "s", "", "Set the secondary model")
	cnfgCmd.Flags().StringVarP(&tertiaryModelFlag, "tertiary", "t", "", "Set the tertiary model")
	cnfgCmd.Flags().Float64VarP(&primaryTempFlag, "primary-temp", "q", 0.0, "Set the primary model temperature (0.1-1.0)")
	cnfgCmd.Flags().Float64VarP(&secondaryTempFlag, "secondary-temp", "w", 0.0, "Set the secondary model temperature (0.1-1.0)")
	cnfgCmd.Flags().Float64VarP(&tertiaryTempFlag, "tertiary-temp", "e", 0.0, "Set the tertiary model temperature (0.1-1.0)")
}
