package configs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lordofthemind/gollama/helpers"
	"github.com/spf13/viper"
)

// LoadGlobalConfig loads the configuration file
func LoadGlobalConfig() (GollamaGlobalConfig, string, error) {
	var config GollamaGlobalConfig

	// Get configuration path
	configPath, err := helpers.GetConfigPath()
	if err != nil {
		return config, "", err
	}

	// Check if the configuration file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("Configuration file not found, creating a new one with default values...")
		config = createDefaultConfig()
		if err := SaveGlobalConfig(config, configPath); err != nil {
			return config, "", err
		}
		fmt.Println("Default configuration file created at:", configPath)
	}

	// Load the configuration using Viper
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return config, "", err
	}

	// Unmarshal the configuration
	if err := viper.Unmarshal(&config); err != nil {
		return config, "", err
	}

	return config, configPath, nil
}

// SaveGlobalConfig saves the configuration file
func SaveGlobalConfig(config GollamaGlobalConfig, configPath string) error {
	viper.Set("primary.model", config.Primary.Model)
	viper.Set("primary.temp", config.Primary.Temp)
	viper.Set("secondary.model", config.Secondary.Model)
	viper.Set("secondary.temp", config.Secondary.Temp)
	viper.Set("tertiary.model", config.Tertiary.Model)
	viper.Set("tertiary.temp", config.Tertiary.Temp)
	viper.Set("setup_completed", config.SetupCompleted)

	// Ensure the config directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), os.ModePerm); err != nil {
		return err
	}

	// Write the configuration to the file
	return viper.WriteConfigAs(configPath)
}

// createDefaultConfig creates a default configuration
func createDefaultConfig() GollamaGlobalConfig {
	return GollamaGlobalConfig{
		Primary: ModelConfig{
			Model: "",
			Temp:  0.0,
		},
		Secondary: ModelConfig{
			Model: "",
			Temp:  0.0,
		},
		Tertiary: ModelConfig{
			Model: "",
			Temp:  0.0,
		},
		SetupCompleted: false,
	}
}
