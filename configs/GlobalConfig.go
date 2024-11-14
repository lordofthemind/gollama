package configs

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

// GollamaGlobalConfig represents the structure of the configuration file
type GollamaGlobalConfig struct {
	PrimaryModel   string  `mapstructure:"primary_model"`
	SecondaryModel string  `mapstructure:"secondary_model"`
	TertiaryModel  string  `mapstructure:"tertiary_model"`
	Temperature    float64 `mapstructure:"temperature"`
	Logging        bool    `mapstructure:"logging"`
	SetupCompleted bool    `mapstructure:"setup_completed"`
}

// LoadConfig loads the configuration file and returns the configuration struct and path
func LoadGlobalConfig() (GollamaGlobalConfig, string, error) {
	var config GollamaGlobalConfig

	// Determine the configuration path based on the OS
	var configPath string
	if customPath := os.Getenv("GOLLAMA_CONFIG"); customPath != "" {
		configPath = customPath
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return config, "", err
		}

		// Set the default config path based on the OS
		if runtime.GOOS == "windows" {
			configPath = filepath.Join(homeDir, "AppData", "Roaming", "Gollama", "gollama.yaml")
		} else {
			configPath = filepath.Join(homeDir, ".config", "Gollama", "gollama.yaml")
		}
	}

	// Check if the config file exists; if not, create a default one
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("Configuration file not found, creating a new one with default values...")

		// Set default values
		config.PrimaryModel = ""
		config.SecondaryModel = ""
		config.TertiaryModel = ""
		config.Temperature = 0.0
		config.SetupCompleted = false

		// Save the default config to file
		if err := SaveGlobalConfig(config, configPath); err != nil {
			return config, "", err
		}

		fmt.Println("Default configuration file created at:", configPath)
	}

	// Load the config file using Viper
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return config, "", err
	}

	// Unmarshal the config into the GollamaGlobalConfig struct
	if err := viper.Unmarshal(&config); err != nil {
		return config, "", err
	}

	return config, configPath, nil
}

// SaveConfig saves the configuration to the specified path
func SaveGlobalConfig(config GollamaGlobalConfig, configPath string) error {
	// Set config values in Viper
	viper.Set("primary_model", config.PrimaryModel)
	viper.Set("secondary_model", config.SecondaryModel)
	viper.Set("tertiary_model", config.TertiaryModel)
	viper.Set("temperature", config.Temperature)
	viper.Set("setup_completed", config.SetupCompleted)

	// Ensure the config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		return err
	}

	// Write the configuration to the file
	return viper.WriteConfigAs(configPath)
}

// DisplayGlobalConfig prints the current configuration from the file
func DisplayGlobalConfig(config GollamaGlobalConfig) {
	fmt.Println("Current Global Configuration:")
	fmt.Printf("Temperature: %.2f\n", config.Temperature)
	fmt.Printf("Primary Model: %s\n", config.PrimaryModel)
	fmt.Printf("Secondary Model: %s\n", config.SecondaryModel)
	fmt.Printf("Tertiary Model: %s\n", config.TertiaryModel)
}
