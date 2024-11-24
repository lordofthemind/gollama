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

// SetupConfiguration runs an interactive setup for configuration
func SetupConfiguration(config *configs.GollamaGlobalConfig, configPath string) {
	models, err := helpers.GetOllamaModels()
	if err != nil {
		fmt.Println("Error retrieving models:", err)
		return
	}

	reader := helpers.NewReader()

	// Select models for Primary, Secondary, and Tertiary
	config.Primary.Model = helpers.SelectModel(reader, models, "Select Primary Model")
	config.Secondary.Model = helpers.SelectModel(reader, models, "Select Secondary Model")
	config.Tertiary.Model = helpers.SelectModel(reader, models, "Select Tertiary Model")

	if helpers.ConfirmAction("Confirm the configuration?") {
		config.SetupCompleted = true
		if err := SaveConfig(*config, configPath); err != nil {
			fmt.Println("Error saving configuration:", err)
		} else {
			fmt.Println("Configuration saved successfully.")
		}
	}
}

// UpdateConfigFromFlags updates the configuration based on CLI flags
func UpdateConfigFromFlags(config *configs.GollamaGlobalConfig, temp float64, pModel, sModel, tModel string) bool {
	models, err := helpers.GetOllamaModels()
	if err != nil {
		fmt.Println("Error retrieving models:", err)
		return false
	}

	updated := false

	if temp != 0.5 {
		config.Primary.Temp = temp
		config.Secondary.Temp = temp
		config.Tertiary.Temp = temp
		updated = true
	}

	if pModel != "" && helpers.ValidateModel(pModel, models) {
		config.Primary.Model = pModel
		updated = true
	}

	if sModel != "" && helpers.ValidateModel(sModel, models) {
		config.Secondary.Model = sModel
		updated = true
	}

	if tModel != "" && helpers.ValidateModel(tModel, models) {
		config.Tertiary.Model = tModel
		updated = true
	}

	return updated
}
