package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/lordofthemind/gollama/configs"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage the application configuration",
	Long: `Configure the application by setting initial values or updating them as needed.
If the configuration is not set, you will be prompted to input the initial setup values.`,
	Run: func(cmd *cobra.Command, args []string) {
		config, configPath, err := configs.LoadGlobalConfig()
		if err != nil {
			fmt.Println("Error loading configuration:", err)
			return
		}

		// Check if setup has already been completed
		if config.SetupCompleted {
			fmt.Println("Configuration already exists:")
			configs.DisplayConfig(config)

			// Ask if the user wants to edit the configuration
			fmt.Print("Do you want to edit the configuration? (yes/no): ")
			if askForConfirmation() {
				config = promptForConfiguration(config)
				fmt.Println("New configuration:")
				configs.DisplayConfig(config)
				if askForConfirmation() {
					if err := configs.SaveGlobalConfig(config, configPath); err != nil {
						fmt.Println("Error saving configuration:", err)
						return
					}
					fmt.Println("Configuration updated successfully.")
				} else {
					fmt.Println("Configuration update canceled.")
				}
			} else {
				fmt.Println("No changes made to the configuration.")
			}
		} else {
			fmt.Println("Initial configuration required.")
			config = promptForConfiguration(config)
			fmt.Println("Please confirm your configuration:")
			configs.DisplayConfig(config)
			if askForConfirmation() {
				config.SetupCompleted = true
				if err := configs.SaveGlobalConfig(config, configPath); err != nil {
					fmt.Println("Error saving configuration:", err)
					return
				}
				fmt.Println("Configuration saved successfully.")
			} else {
				fmt.Println("Configuration setup canceled.")
			}
		}
	},
}

// init function adds the configCmd to rootCmd
func init() {
	rootCmd.AddCommand(configCmd)
}

// promptForConfiguration prompts the user for each configuration setting
func promptForConfiguration(config configs.GollamaGlobalConfig) configs.GollamaGlobalConfig {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Ollama URL: ")
	config.OllamaURL = readInput(reader)

	fmt.Print("Enter Primary Model: ")
	config.PrimaryModel = readInput(reader)

	fmt.Print("Enter Secondary Model: ")
	config.SecondaryModel = readInput(reader)

	fmt.Print("Enter Tertiary Model: ")
	config.TertiaryModel = readInput(reader)

	fmt.Print("Enter Temperature (e.g., 0.5): ")
	config.Temperature = readFloatInput(reader)

	return config
}

// readInput reads input from the user and trims whitespace
func readInput(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// readFloatInput reads a float input from the user and returns it
func readFloatInput(reader *bufio.Reader) float64 {
	var value float64
	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		_, err := fmt.Sscanf(input, "%f", &value)
		if err == nil {
			break
		}
		fmt.Print("Invalid input. Please enter a valid number: ")
	}
	return value
}

// askForConfirmation prompts the user to confirm an action
func askForConfirmation() bool {
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "yes" || response == "y"
}
