package services

import (
	"fmt"

	"github.com/lordofthemind/gollama/configs"
	"github.com/lordofthemind/gollama/helpers"
)

// LoadConfig loads the configuration file
func LoadConfig() (configs.GollamaGlobalConfig, string, error) {
	return configs.LoadGlobalConfig()
}

// SaveConfig saves the configuration file
func SaveConfig(config configs.GollamaGlobalConfig, path string) error {
	return configs.SaveGlobalConfig(config, path)
}

// DisplayConfig displays the current configuration
func DisplayConfig(config configs.GollamaGlobalConfig) {
	configs.DisplayGlobalConfig(config)
}

// SelectModelWithTemperature prompts the user to select models and input temperatures
func SelectModelWithTemperature(config *configs.GollamaGlobalConfig) {
	// Retrieve available models
	models, err := helpers.GetOllamaModels()
	if err != nil {
		fmt.Println("Error retrieving models:", err)
		return
	}

	// Ensure there are available models
	if len(models) == 0 {
		fmt.Println("No models available. Please pull models using `ollama pull <model_name>` and try again.")
		return
	}

	// Select Primary Model
	selectedPrimary := helpers.SelectModel(nil, models, "Select Primary Model")
	config.Primary.Model = selectedPrimary
	config.Primary.Temp = helpers.PromptForTemperature("Enter temperature for Primary Model (0.1-1.0): ")

	// Select Secondary Model
	selectedSecondary := helpers.SelectModel(nil, models, "Select Secondary Model")
	config.Secondary.Model = selectedSecondary
	config.Secondary.Temp = helpers.PromptForTemperature("Enter temperature for Secondary Model (0.1-1.0): ")

	// Select Tertiary Model
	selectedTertiary := helpers.SelectModel(nil, models, "Select Tertiary Model")
	config.Tertiary.Model = selectedTertiary
	config.Tertiary.Temp = helpers.PromptForTemperature("Enter temperature for Tertiary Model (0.1-1.0): ")
}

// UpdateConfigFromFlags updates the configuration based on CLI flags
func UpdateConfigFromFlags(
	config *configs.GollamaGlobalConfig,
	primaryTemp, secondaryTemp, tertiaryTemp float64,
	pModel, sModel, tModel string,
) bool {
	models, err := helpers.GetOllamaModels()
	if err != nil {
		fmt.Println("Error retrieving models:", err)
		return false
	}

	updated := false

	// Update primary temperature
	if primaryTemp != 0.0 {
		if err := helpers.ValidateTemperature(primaryTemp); err != nil {
			fmt.Println("Error:", err)
			return false
		}
		config.Primary.Temp = primaryTemp
		updated = true
	}

	// Update secondary temperature
	if secondaryTemp != 0.0 {
		if err := helpers.ValidateTemperature(secondaryTemp); err != nil {
			fmt.Println("Error:", err)
			return false
		}
		config.Secondary.Temp = secondaryTemp
		updated = true
	}

	// Update tertiary temperature
	if tertiaryTemp != 0.0 {
		if err := helpers.ValidateTemperature(tertiaryTemp); err != nil {
			fmt.Println("Error:", err)
			return false
		}
		config.Tertiary.Temp = tertiaryTemp
		updated = true
	}

	// Update primary model
	if pModel != "" && helpers.ValidateModel(pModel, models) {
		config.Primary.Model = pModel
		updated = true
	}

	// Update secondary model
	if sModel != "" && helpers.ValidateModel(sModel, models) {
		config.Secondary.Model = sModel
		updated = true
	}

	// Update tertiary model
	if tModel != "" && helpers.ValidateModel(tModel, models) {
		config.Tertiary.Model = tModel
		updated = true
	}

	return updated
}
