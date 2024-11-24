package configs

// import (
// 	"fmt"
// 	"os"
// 	"path/filepath"
// 	"runtime"

// 	"github.com/spf13/viper"
// )

// // GollamaGlobalConfig represents the updated structure of the configuration file
// type GollamaGlobalConfig struct {
// 	Primary struct {
// 		Model string  `mapstructure:"model"`
// 		Temp  float64 `mapstructure:"temp"`
// 	} `mapstructure:"primary"`
// 	Secondary struct {
// 		Model string  `mapstructure:"model"`
// 		Temp  float64 `mapstructure:"temp"`
// 	} `mapstructure:"secondary"`
// 	Tertiary struct {
// 		Model string  `mapstructure:"model"`
// 		Temp  float64 `mapstructure:"temp"`
// 	} `mapstructure:"tertiary"`
// 	SetupCompleted bool `mapstructure:"setup_completed"`
// }

// // LoadGlobalConfig loads the configuration file and returns the configuration struct and path
// func LoadGlobalConfig() (GollamaGlobalConfig, string, error) {
// 	var config GollamaGlobalConfig

// 	// Determine the configuration path based on the OS
// 	var configPath string
// 	if customPath := os.Getenv("GOLLAMA_CONFIG"); customPath != "" {
// 		configPath = customPath
// 	} else {
// 		homeDir, err := os.UserHomeDir()
// 		if err != nil {
// 			return config, "", err
// 		}

// 		// Set the default config path based on the OS
// 		if runtime.GOOS == "windows" {
// 			configPath = filepath.Join(homeDir, "AppData", "Roaming", "Gollama", "gollama.yaml")
// 		} else {
// 			configPath = filepath.Join(homeDir, ".config", "Gollama", "gollama.yaml")
// 		}
// 	}

// 	// Check if the config file exists; if not, create a default one
// 	if _, err := os.Stat(configPath); os.IsNotExist(err) {
// 		fmt.Println("Configuration file not found, creating a new one with default values...")

// 		// Set default values
// 		config = GollamaGlobalConfig{
// 			Primary: struct {
// 				Model string  `mapstructure:"model"`
// 				Temp  float64 `mapstructure:"temp"`
// 			}{Model: "", Temp: 0.0},
// 			Secondary: struct {
// 				Model string  `mapstructure:"model"`
// 				Temp  float64 `mapstructure:"temp"`
// 			}{Model: "", Temp: 0.0},
// 			Tertiary: struct {
// 				Model string  `mapstructure:"model"`
// 				Temp  float64 `mapstructure:"temp"`
// 			}{Model: "", Temp: 0.0},
// 			SetupCompleted: false,
// 		}

// 		// Save the default config to file
// 		if err := SaveGlobalConfig(config, configPath); err != nil {
// 			return config, "", err
// 		}

// 		fmt.Println("Default configuration file created at:", configPath)
// 	}

// 	// Load the config file using Viper
// 	viper.SetConfigFile(configPath)
// 	if err := viper.ReadInConfig(); err != nil {
// 		return config, "", err
// 	}

// 	// Unmarshal the config into the GollamaGlobalConfig struct
// 	if err := viper.Unmarshal(&config); err != nil {
// 		return config, "", err
// 	}

// 	return config, configPath, nil
// }

// // SaveGlobalConfig saves the configuration to the specified path
// func SaveGlobalConfig(config GollamaGlobalConfig, configPath string) error {
// 	// Set nested config values in Viper
// 	viper.Set("primary.model", config.Primary.Model)
// 	viper.Set("primary.temp", config.Primary.Temp)
// 	viper.Set("secondary.model", config.Secondary.Model)
// 	viper.Set("secondary.temp", config.Secondary.Temp)
// 	viper.Set("tertiary.model", config.Tertiary.Model)
// 	viper.Set("tertiary.temp", config.Tertiary.Temp)
// 	viper.Set("setup_completed", config.SetupCompleted)

// 	// Ensure the config directory exists
// 	configDir := filepath.Dir(configPath)
// 	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
// 		return err
// 	}

// 	// Write the configuration to the file
// 	return viper.WriteConfigAs(configPath)
// }

// // DisplayGlobalConfig prints the current configuration from the file
// func DisplayGlobalConfig(config GollamaGlobalConfig) {
// 	fmt.Println("Current Global Configuration:")
// 	fmt.Println("Primary:")
// 	fmt.Printf("  Model: %s\n", config.Primary.Model)
// 	fmt.Printf("  Temperature: %.2f\n", config.Primary.Temp)
// 	fmt.Println("Secondary:")
// 	fmt.Printf("  Model: %s\n", config.Secondary.Model)
// 	fmt.Printf("  Temperature: %.2f\n", config.Secondary.Temp)
// 	fmt.Println("Tertiary:")
// 	fmt.Printf("  Model: %s\n", config.Tertiary.Model)
// 	fmt.Printf("  Temperature: %.2f\n", config.Tertiary.Temp)
// 	fmt.Printf("Setup Completed: %t\n", config.SetupCompleted)
// }
