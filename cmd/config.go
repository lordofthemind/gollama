package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/lordofthemind/gollama/configs"
)

// Function to check if `ollama` is installed
func checkOllamaInstallation() bool {
	_, err := exec.LookPath("ollama")
	if err != nil {
		fmt.Println("Ollama is not installed. Please install Ollama to proceed with Gollama setup.")
		return false
	}
	return true
}

// Function to retrieve models available in `ollama`
func getOllamaModels() ([]string, error) {
	cmd := exec.Command("ollama", "list")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	models := []string{}
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] != "NAME" { // Ignore header line
			models = append(models, fields[0])
		}
	}
	return models, nil
}

// Function to initiate first-time setup
func initiateSetup(config *configs.GollamaGlobalConfig, configPath string) {
	// Ensure Ollama is installed before proceeding
	if !checkOllamaInstallation() {
		return
	}

	// Retrieve list of available models
	models, err := getOllamaModels()
	if err != nil {
		fmt.Println("Error retrieving models from Ollama:", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Initial setup required. Please provide the following details:")

	// Prompt for Ollama URL
	fmt.Print("Enter Ollama URL: ")
	config.OllamaURL = readInput(reader)

	// Prompt for Primary Model with validation from list of available models
	fmt.Println("Available models:")
	for i, model := range models {
		fmt.Printf("%d. %s\n", i+1, model)
	}
	fmt.Print("Enter the number corresponding to your Primary Model: ")
	config.PrimaryModel = selectModel(reader, models)

	// Prompt for Secondary Model
	fmt.Print("Enter the number corresponding to your Secondary Model: ")
	config.SecondaryModel = selectModel(reader, models)

	// Prompt for Tertiary Model
	fmt.Print("Enter the number corresponding to your Tertiary Model: ")
	config.TertiaryModel = selectModel(reader, models)

	// Prompt for Temperature
	config.Temperature = readFloatInput(reader)

	// Confirmation prompt
	for {
		fmt.Println("\nPlease confirm the details you entered:")
		configs.DisplayConfig(*config)
		fmt.Print("Are these details correct? (y/n): ")

		confirmation := readInput(reader)
		if confirmation == "y" {
			// Set SetupCompleted to true and save config
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
			initiateSetup(config, configPath) // Recursive call to re-enter setup
			break
		} else {
			fmt.Println("Invalid option. Please type 'y' for yes or 'n' for no.")
		}
	}
}

// Helper function to select a model based on user input
func selectModel(reader *bufio.Reader, models []string) string {
	for {
		modelIndex, _ := strconv.Atoi(readInput(reader))
		if modelIndex > 0 && modelIndex <= len(models) {
			return models[modelIndex-1]
		} else {
			fmt.Printf("Invalid choice. Please select a number between 1 and %d: ", len(models))
		}
	}
}

// readInput reads input from the user and trims whitespace
func readInput(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// readFloatInput reads a float input from the user
func readFloatInput(reader *bufio.Reader) float64 {
	var value float64
	for {
		fmt.Print("Enter Temperature (e.g., 0.5): ")
		input := readInput(reader)
		_, err := fmt.Sscanf(input, "%f", &value)
		if err == nil {
			break
		}
		fmt.Println("Invalid input. Please enter a valid number.")
	}
	return value
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage the application configuration",
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

			fmt.Print("Do you want to edit the configuration? (yes/no): ")
			if readInput(bufio.NewReader(os.Stdin)) == "yes" {
				initiateSetup(&config, configPath)
			} else {
				fmt.Println("No changes made to the configuration.")
			}
		} else {
			initiateSetup(&config, configPath)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
