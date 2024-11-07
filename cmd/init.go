package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/lordofthemind/gollama/configs"
	"github.com/lordofthemind/gollama/utils"
)

var projectNameFlag string

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Gollama workspace configuration",
	Long: `The init command sets up a new Gollama configuration in the current working directory.
If no project name is specified, it defaults to the parent directory name or prompts the user.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Determine the project name
		projectName := projectNameFlag
		if projectName == "" {
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Println("Error getting current directory:", err)
				return
			}
			projectName = filepath.Base(cwd)
			if projectName == "" {
				fmt.Print("Enter the project name: ")
				projectName = utils.ReadInput(bufio.NewReader(os.Stdin))
			}
		}

		// Set up configuration path in the current working directory
		configDir := filepath.Join(".", ".gollama")
		configPath := filepath.Join(configDir, "gollama.yaml")

		// Check if the configuration directory exists, if not, create it
		if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
			fmt.Println("Error creating configuration directory:", err)
			return
		}

		// Initialize and configure GollamaWorkspaceConfig
		var config configs.GollamaWorkspaceConfig
		config.ProjectName = projectName
		initiateWorkspaceConfigurationSetup(&config, configPath)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&projectNameFlag, "project", "p", "", "Project name for the new Gollama workspace")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// Main setup function for the workspace configuration
func initiateWorkspaceConfigurationSetup(config *configs.GollamaWorkspaceConfig, configPath string) {
	if !utils.CheckOllamaInstallation() {
		return
	}

	models, err := utils.GetOllamaModels()
	if err != nil {
		fmt.Println("Error retrieving models from Ollama:", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Available models:")
	for i, model := range models {
		fmt.Printf("%d. %s\n", i+1, model)
	}

	// Select models for configuration
	fmt.Print("Select Primary Model by entering the corresponding number: ")
	config.PrimaryModel = utils.SelectModel(reader, models)
	fmt.Print("Select Secondary Model by entering the corresponding number: ")
	config.SecondaryModel = utils.SelectModel(reader, models)
	fmt.Print("Select Tertiary Model by entering the corresponding number: ")
	config.TertiaryModel = utils.SelectModel(reader, models)

	// Set temperature
	utils.DisplayTemperatureGuidance()
	config.Temperature = utils.ReadFloatInput(reader)

	// Confirmation prompt
	for {
		fmt.Println("\nPlease confirm the details you entered:")
		configs.DisplayWorkspaceConfig(*config)
		fmt.Print("Are these details correct? (y/n): ")

		confirmation := utils.ReadInput(reader)
		if confirmation == "y" {
			err := configs.SaveWorkspaceConfig(*config, configPath)
			if err != nil {
				fmt.Println("Error saving configuration:", err)
				return
			}
			fmt.Println("Configuration saved successfully.")
			break
		} else if confirmation == "n" {
			fmt.Println("Let's re-enter the details.")
			initiateWorkspaceConfigurationSetup(config, configPath)
			break
		} else {
			fmt.Println("Invalid option. Please type 'y' for yes or 'n' for no.")
		}
	}
}
