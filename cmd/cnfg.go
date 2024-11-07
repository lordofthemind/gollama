/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/lordofthemind/gollama/configs"
	"github.com/lordofthemind/gollama/utils"
)

const defaultOllamaURL = "http://localhost:11434/"

var (
	tempFlag   float64
	pModelFlag string
	sModelFlag string
	tModelFlag string
)

// cnfgCmd represents the cnfg command
var cnfgCmd = &cobra.Command{
	Use:   "cnfg",
	Short: "Manage the application configuration",
	Long:  "The config command helps set up and manage Gollama’s configuration...",
	Run: func(cmd *cobra.Command, args []string) {
		config, configPath, err := configs.LoadGlobalConfig()
		if err != nil {
			fmt.Println("Error loading configuration:", err)
			return
		}

		// Check if config file exists; if not, create it with full setup
		if !config.SetupCompleted {
			initiateSetup(&config, configPath)
			return
		}

		// If no flags, show current configuration and prompt for edit
		if !cmd.Flags().Changed("temp") && !cmd.Flags().Changed("pmodel") &&
			!cmd.Flags().Changed("smodel") && !cmd.Flags().Changed("tmodel") {
			fmt.Println("Configuration already exists:")
			configs.DisplayGlobalConfig(config)
			fmt.Print("Do you want to edit the configuration? (y/n): ")
			if utils.ReadInput(bufio.NewReader(os.Stdin)) == "y" {
				initiateSetup(&config, configPath)
			} else {
				fmt.Println("No changes made to the configuration.")
			}
			return
		}

		// Track updates
		anyUpdate := false
		models, err := utils.GetOllamaModels()
		if err != nil {
			fmt.Println("Error retrieving models:", err)
			return
		}

		// Apply flag-based updates
		if tempFlag != 0.5 {
			config.Temperature = tempFlag
			fmt.Printf("Temperature updated to: %.2f\n", config.Temperature)
			anyUpdate = true
		}
		if pModelFlag != "" {
			if utils.ValidateModel(pModelFlag, models) {
				config.PrimaryModel = pModelFlag
				fmt.Printf("Primary Model updated to: %s\n", config.PrimaryModel)
				anyUpdate = true
			} else {
				return
			}
		}
		if sModelFlag != "" {
			if utils.ValidateModel(sModelFlag, models) {
				config.SecondaryModel = sModelFlag
				fmt.Printf("Secondary Model updated to: %s\n", config.SecondaryModel)
				anyUpdate = true
			} else {
				return
			}
		}
		if tModelFlag != "" {
			if utils.ValidateModel(tModelFlag, models) {
				config.TertiaryModel = tModelFlag
				fmt.Printf("Tertiary Model updated to: %s\n", config.TertiaryModel)
				anyUpdate = true
			} else {
				return
			}
		}

		// Save only if updates were made
		if anyUpdate {
			err = configs.SaveGlobalConfig(config, configPath)
			if err != nil {
				fmt.Println("Error saving configuration:", err)
			} else {
				fmt.Println("Configuration updated successfully.")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(cnfgCmd)
	cnfgCmd.Flags().StringVarP(&pModelFlag, "primary", "p", "", "Primary model name")
	cnfgCmd.Flags().StringVarP(&sModelFlag, "secondary", "s", "", "Secondary model name")
	cnfgCmd.Flags().StringVarP(&tModelFlag, "tertiary", "e", "", "Tertiary model name")
	cnfgCmd.Flags().Float64VarP(&tempFlag, "temp", "t", 0.5, "Temperature setting for the model (0.1 to 1.0)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cnfgCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cnfgCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// Main setup function for initial or reconfiguration
func initiateSetup(config *configs.GollamaGlobalConfig, configPath string) {
	if !utils.CheckOllamaInstallation() {
		return
	}

	models, err := utils.GetOllamaModels()
	if err != nil {
		fmt.Println("Error retrieving models from Ollama:", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	if !config.SetupCompleted {
		fmt.Println("Initial setup required. Please provide the following details:")
	}
	fmt.Printf("Enter Ollama URL (default: %s): ", defaultOllamaURL)
	config.OllamaURL = utils.ReadInput(reader)
	if config.OllamaURL == "" {
		config.OllamaURL = defaultOllamaURL
	}

	fmt.Println("Available models:")
	for i, model := range models {
		fmt.Printf("%d. %s\n", i+1, model)
	}

	fmt.Print("Select Primary Model by entering the corresponding number: ")
	config.PrimaryModel = utils.SelectModel(reader, models)
	fmt.Print("Select Secondary Model by entering the corresponding number: ")
	config.SecondaryModel = utils.SelectModel(reader, models)
	fmt.Print("Select Tertiary Model by entering the corresponding number: ")
	config.TertiaryModel = utils.SelectModel(reader, models)

	utils.DisplayTemperatureGuidance()
	config.Temperature = utils.ReadFloatInput(reader)

	// Confirmation prompt
	for {
		fmt.Println("\nPlease confirm the details you entered:")
		configs.DisplayGlobalConfig(*config)
		fmt.Print("Are these details correct? (y/n): ")

		confirmation := utils.ReadInput(reader)
		if confirmation == "y" {
			config.SetupCompleted = true
			err := configs.SaveGlobalConfig(*config, configPath)
			if err != nil {
				fmt.Println("Error saving configuration:", err)
				return
			}
			fmt.Println("Configuration saved successfully.")
			break
		} else if confirmation == "n" {
			fmt.Println("Let's re-enter the details.")
			initiateSetup(config, configPath)
			break
		} else {
			fmt.Println("Invalid option. Please type 'y' for yes or 'n' for no.")
		}
	}
}
