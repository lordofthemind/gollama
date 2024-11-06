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

const defaultOllamaURL = "http://localhost:11434/" // Default Ollama URL

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

// Main setup function for initial or reconfiguration
func initiateSetup(config *configs.GollamaGlobalConfig, configPath string) {
	if !checkOllamaInstallation() {
		return
	}

	models, err := getOllamaModels()
	if err != nil {
		fmt.Println("Error retrieving models from Ollama:", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)
	if !config.SetupCompleted {
		fmt.Println("Initial setup required. Please provide the following details:")
	}
	fmt.Printf("Enter Ollama URL (default: %s): ", defaultOllamaURL)
	config.OllamaURL = readInput(reader)
	if config.OllamaURL == "" {
		config.OllamaURL = defaultOllamaURL
	}

	fmt.Println("Available models:")
	for i, model := range models {
		fmt.Printf("%d. %s\n", i+1, model)
	}

	fmt.Print("Select Primary Model by entering the corresponding number: ")
	config.PrimaryModel = selectModel(reader, models)
	fmt.Print("Select Secondary Model by entering the corresponding number: ")
	config.SecondaryModel = selectModel(reader, models)
	fmt.Print("Select Tertiary Model by entering the corresponding number: ")
	config.TertiaryModel = selectModel(reader, models)

	displayTemperatureGuidance()
	config.Temperature = readFloatInput(reader)

	// Confirmation prompt
	for {
		fmt.Println("\nPlease confirm the details you entered:")
		configs.DisplayConfig(*config)
		fmt.Print("Are these details correct? (y/n): ")

		confirmation := readInput(reader)
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

// Helper function to display temperature guidance
func displayTemperatureGuidance() {
	fmt.Println("\n### Temperature Guidance ###")
	fmt.Println("Temperature settings guide:")
	fmt.Println("0.0 - 0.3: Deterministic, ideal for precise tasks")
	fmt.Println("0.4 - 0.7: Balanced, suitable for conversations")
	fmt.Println("0.8 - 1.0: High randomness, for creative tasks")
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
		fmt.Print("Enter Temperature (0.1 - 1.0, e.g., 0.5): ")
		input := readInput(reader)
		_, err := fmt.Sscanf(input, "%f", &value)
		if err == nil && value >= 0.1 && value <= 1.0 {
			break
		}
		fmt.Println("Invalid input. Please enter a valid number between 0.1 and 1.0.")
	}
	return value
}

var (
	tempFlag   float64
	pModelFlag string
	sModelFlag string
	tModelFlag string
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage the application configuration",
	Long:  "",
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
			configs.DisplayConfig(config)
			fmt.Print("Do you want to edit the configuration? (y/n): ")
			if readInput(bufio.NewReader(os.Stdin)) == "y" {
				initiateSetup(&config, configPath)
			} else {
				fmt.Println("No changes made to the configuration.")
			}
			return
		}

		// Track updates
		anyUpdate := false
		models, err := getOllamaModels()
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
			if validateModel(pModelFlag, models) {
				config.PrimaryModel = pModelFlag
				fmt.Printf("Primary Model updated to: %s\n", config.PrimaryModel)
				anyUpdate = true
			} else {
				return
			}
		}
		if sModelFlag != "" {
			if validateModel(sModelFlag, models) {
				config.SecondaryModel = sModelFlag
				fmt.Printf("Secondary Model updated to: %s\n", config.SecondaryModel)
				anyUpdate = true
			} else {
				return
			}
		}
		if tModelFlag != "" {
			if validateModel(tModelFlag, models) {
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

func validateModel(model string, models []string) bool {
	for _, m := range models {
		if m == model {
			return true
		}
	}
	fmt.Printf("Model %s is not available in Ollama installation.\n", model)
	return false
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().Float64VarP(&tempFlag, "temp", "t", 0.5, "Temperature setting for the model (0.1 to 1.0)")
	configCmd.Flags().StringVarP(&pModelFlag, "pmodel", "p", "", "Primary model name")
	configCmd.Flags().StringVarP(&sModelFlag, "smodel", "s", "", "Secondary model name")
	configCmd.Flags().StringVarP(&tModelFlag, "tmodel", "e", "", "Tertiary model name")
}
