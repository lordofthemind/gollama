package configs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// GollamaWorkspaceConfig represents the structure of the configuration file
type GollamaWorkspaceConfig struct {
	ProjectName    string  `mapstructure:"project_name"`
	PrimaryModel   string  `mapstructure:"primary_model"`
	SecondaryModel string  `mapstructure:"secondary_model"`
	TertiaryModel  string  `mapstructure:"tertiary_model"`
	Temperature    float64 `mapstructure:"temperature"`
}

func SaveWorkspaceConfig(config GollamaWorkspaceConfig, configPath string) error {
	viper.Set("project_name", config.ProjectName)
	viper.Set("primary_model", config.PrimaryModel)
	viper.Set("secondary_model", config.SecondaryModel)
	viper.Set("tertiary_model", config.TertiaryModel)
	viper.Set("temperature", config.Temperature)

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		return err
	}

	return viper.WriteConfigAs(configPath)
}

// LoadWorkspaceConfig loads the workspace configuration file and returns the configuration struct and path
func LoadWorkspaceConfig() (GollamaWorkspaceConfig, string, error) {
	var config GollamaWorkspaceConfig

	// Define the workspace configuration path based on the current directory
	cwd, err := os.Getwd()
	if err != nil {
		return config, "", err
	}
	configPath := filepath.Join(cwd, ".gollama", "gollama.yaml")

	// Check if the workspace config file exists; if not, return an error
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, "", fmt.Errorf("configuration file not found at: %s", configPath)
	}

	// Load the config file using Viper
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return config, "", err
	}

	// Unmarshal the config into the GollamaWorkspaceConfig struct
	if err := viper.Unmarshal(&config); err != nil {
		return config, "", err
	}

	return config, configPath, nil
}

func DisplayWorkspaceConfig(config GollamaWorkspaceConfig) {
	fmt.Println("Current Workspace Configuration:")
	fmt.Printf("Project Name: %s\n", config.ProjectName)
	fmt.Printf("Primary Model: %s\n", config.PrimaryModel)
	fmt.Printf("Secondary Model: %s\n", config.SecondaryModel)
	fmt.Printf("Tertiary Model: %s\n", config.TertiaryModel)
	fmt.Printf("Temperature: %.2f\n", config.Temperature)
}
