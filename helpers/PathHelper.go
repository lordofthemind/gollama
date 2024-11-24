package helpers

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetConfigPath determines the path of the configuration file
func GetConfigPath() (string, error) {
	if customPath := os.Getenv("GOLLAMA_CONFIG"); customPath != "" {
		return customPath, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	if runtime.GOOS == "windows" {
		return filepath.Join(homeDir, "AppData", "Roaming", "Gollama", "gollama.yaml"), nil
	}
	return filepath.Join(homeDir, ".config", "Gollama", "gollama.yaml"), nil
}
