package configs

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

// GollamaConfig represents the combined configuration structure
type GollamaConfig struct {
	ProjectName    string  `mapstructure:"project_name"`
	OllamaURL      string  `mapstructure:"ollama_url"`
	PrimaryModel   string  `mapstructure:"primary_model"`
	SecondaryModel string  `mapstructure:"secondary_model"`
	TertiaryModel  string  `mapstructure:"tertiary_model"`
	Temperature    float64 `mapstructure:"temperature"`
}

// LoadConfig loads both workspace and global configuration, prioritizing workspace values if available.
func LoadConfig() (GollamaConfig, error) {
	var config GollamaConfig

	// Load workspace config first (if it exists)
	cwd, err := os.Getwd()
	if err != nil {
		return config, err
	}
	workspaceConfigPath := filepath.Join(cwd, ".gollama", "gollama.yaml")

	// If workspace config exists, load it
	if _, err := os.Stat(workspaceConfigPath); err == nil {
		viper.SetConfigFile(workspaceConfigPath)
		if err := viper.ReadInConfig(); err != nil {
			return config, err
		}
		if err := viper.Unmarshal(&config); err != nil {
			return config, err
		}
	}

	// Load global config and override values only if not already set by workspace config
	globalConfig, _, err := loadGlobalConfigOnly()
	if err != nil {
		return config, err
	}

	// Override any unset fields in workspace config with values from global config
	if config.ProjectName == "" {
		config.ProjectName = globalConfig.ProjectName
	}
	if config.OllamaURL == "" {
		config.OllamaURL = globalConfig.OllamaURL
	}
	if config.PrimaryModel == "" {
		config.PrimaryModel = globalConfig.PrimaryModel
	}
	if config.SecondaryModel == "" {
		config.SecondaryModel = globalConfig.SecondaryModel
	}
	if config.TertiaryModel == "" {
		config.TertiaryModel = globalConfig.TertiaryModel
	}
	if config.Temperature == 0.0 {
		config.Temperature = globalConfig.Temperature
	}

	return config, nil
}

// loadGlobalConfigOnly loads only the global configuration file.
func loadGlobalConfigOnly() (GollamaConfig, string, error) {
	var config GollamaConfig

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
			configPath = filepath.Join(homeDir, "AppData", "Roaming", "gollama", "gollama.yaml") // Windows
		} else {
			configPath = filepath.Join(homeDir, ".config", "gollama", "gollama.yaml") // Linux/Mac
		}
	}

	// Load the global config file using Viper
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return config, "", err
	}

	// Unmarshal the config into the GollamaConfig struct
	if err := viper.Unmarshal(&config); err != nil {
		return config, "", err
	}

	return config, configPath, nil
}
